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
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type Event interface {
	GetID() string
	GetType() string
	GetTimestamp() time.Time
	GetPayload() interface{}
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

type EventHandler interface {
	Handle(ctx context.Context, event Event) error
	CanHandle(eventType string) bool
}

type EventBus struct {
	handlers map[string][]EventHandler
	mutex    sync.RWMutex
	logger   *log.Logger
}

//func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
//	eb.mutex.Lock()
//	defer eb.mutex.Unlock()
//
//	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
//	eb.logger.Printf("Subscribed handler for event type: %s", eventType)
//}

func (eb *EventBus) Publish(ctx context.Context, event Event) error {
	eb.logger.Printf("Publishing event: %s (ID: %s)", event.GetType(), event.GetID())

	payloadInBytes, err := json.Marshal(event.GetPayload())
	if err != nil {
		return err
	}
	err = pushMessageToQueue(event.GetType(), payloadInBytes)
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

func NewEventBus(logger *log.Logger) *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
		logger:   logger,
	}
}

func connectProducer() (sarama.SyncProducer, error) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	log.Printf("loading producer config: %s", cfg.KafkaBootstrapServers)
	return sarama.NewSyncProducer([]string{cfg.KafkaBootstrapServers}, cfg.GetSeranaConfig())
}

func pushMessageToQueue(topic string, message []byte) error {

	producer, err := connectProducer()
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

func ConnectConsumer() (sarama.Consumer, error) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	return sarama.NewConsumer([]string{cfg.KafkaBootstrapServers}, cfg.GetSeranaConfig())
}

func StartConsumers(ctx context.Context, handlers map[string]EventHandler) {
	var consumers []sarama.PartitionConsumer
	var workers []sarama.Consumer

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for topic, handler := range handlers {
		worker, err := ConnectConsumer()
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
					concreteEvent, err := CreateEventFromType(topic)
					if err != nil {
						fmt.Printf("Error creating event: %v\n", err)
						continue
					}

					if err := json.Unmarshal(msg.Value, &concreteEvent); err != nil {
						fmt.Printf("Error unmarshaling event: %v\n", err)
						continue
					}

					if handler.CanHandle(topic) {
						err := handler.Handle(ctx, concreteEvent)
						if err != nil {
							fmt.Printf("Error handling message: %v\n", err)
							continue
						}
					}
				}
			}
		}(topic, handler, consumer)
	}

	// Wait for termination signal
	<-sigchan

	// Clean up all consumers
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

type EventFactory func() Event

var eventRegistry = make(map[string]EventFactory)

func RegisterEventType(eventType string, factory EventFactory) {
	eventRegistry[eventType] = factory
}

func CreateEventFromType(eventType string) (Event, error) {
	factory, exists := eventRegistry[eventType]
	if !exists {
		return nil, fmt.Errorf("no factory registered for event type: %s", eventType)
	}
	return factory(), nil
}
