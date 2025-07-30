package shared

import (
	"context"
	"encoding/json"
	"event-driven-go/internal/config"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
)

var GlobalEventBus *EventBus

func init() {
	GlobalEventBus = NewEventBus()
}

type Event interface {
	GetID() string
	GetType() string
	GetTimestamp() time.Time
	GetPayload() interface{}
}

type EventHandler interface {
	Handle(ctx context.Context, event Event) error
	CanHandle(eventType string) bool
}

type BaseEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

func (e BaseEvent) GetID() string {
	return e.ID
}

func (e BaseEvent) GetType() string {
	return e.Type
}

func (e BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

type EventBus struct {
	eventRegistry map[string]EventRegistration
}

type EventRegistration struct {
	eventHandler EventHandler
	event        Event
}

func (eb EventBus) RegisterEventType(eventType string, event Event, handler EventHandler) {
	eb.eventRegistry[eventType] = EventRegistration{
		handler,
		event,
	}
}

func (eb *EventBus) Publish(ctx context.Context, event Event) error {
	log.Printf("Publishing event: %s (ID: %s)", event.GetType(), event.GetID())

	payloadInBytes, err := json.Marshal(event.GetPayload())
	if err != nil {
		return err
	}
	err = eb.pushMessageToQueue(event.GetType(), payloadInBytes)
	if err != nil {
		return err
	}

	return nil
}

func NewBaseEvent(eventType string) BaseEvent {
	return BaseEvent{
		ID:        uuid.New().String(),
		Type:      eventType,
		Timestamp: time.Now(),
	}
}

func NewEventBus() *EventBus {
	return &EventBus{
		eventRegistry: make(map[string]EventRegistration),
	}
}

func (eb EventBus) connectProducer() (sarama.SyncProducer, error) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	log.Printf("loading producer config: %s", cfg.KafkaBootstrapServers)
	return sarama.NewSyncProducer([]string{cfg.KafkaBootstrapServers}, cfg.GetSeranaConfig())
}

func (eb EventBus) pushMessageToQueue(topic string, message []byte) error {
	producer, err := eb.connectProducer()
	if err != nil {
		return err
	}
	defer producer.Close()

	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(producerMessage)
	if err != nil {
		return err
	}

	fmt.Printf("Topic: %s, Partition: %d, Offset: %d\n", topic, partition, offset)

	return nil
}

func (eb EventBus) connectConsumer() (sarama.Consumer, error) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	return sarama.NewConsumer([]string{cfg.KafkaBootstrapServers}, cfg.GetSeranaConfig())
}

func (eb EventBus) StartConsumers(ctx context.Context) {
	var consumers []sarama.PartitionConsumer
	var workers []sarama.Consumer

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for topic, registration := range eb.eventRegistry {
		worker, err := eb.connectConsumer()
		if err != nil {
			panic(err)
		}
		workers = append(workers, worker)

		consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
		if err != nil {
			panic(err)
		}
		consumers = append(consumers, consumer)

		fmt.Printf("Consuming messages for %s\n", topic)

		go func(topic string, handler EventHandler, consumer sarama.PartitionConsumer) {
			for {
				select {
				case err := <-consumer.Errors():
					fmt.Printf("Error consuming message: %v\n", err)
				case msg := <-consumer.Messages():
					if err != nil {
						fmt.Printf("Error creating event: %v\n", err)
						continue
					}

					if err := json.Unmarshal(msg.Value, &registration.event); err != nil {
						fmt.Printf("Error unmarshaling event: %v\n", err)
						continue
					}

					if handler.CanHandle(topic) {
						err := handler.Handle(ctx, registration.event)
						if err != nil {
							fmt.Printf("Error handling message: %v\n", err)
							continue
						}
					}
				}
			}
		}(topic, registration.eventHandler, consumer)
	}

	<-sigchan

	for _, consumer := range consumers {
		if err := consumer.Close(); err != nil {
			fmt.Printf("Error closing consumer: %v\n", err)
		}
	}

	for _, worker := range workers {
		if err := worker.Close(); err != nil {
			fmt.Printf("Error closing worker: %v\n", err)
		}
	}
}
