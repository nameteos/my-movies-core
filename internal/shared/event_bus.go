package shared

import (
	"context"
	"encoding/json"
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
	SyncProducer  sarama.SyncProducer
	Consumer      sarama.Consumer
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
	syncProducer, err := sarama.NewSyncProducer([]string{Config.Kafka.BootstrapServers}, Config.GetSaramaConfig())
	if err != nil {
		log.Fatal(err)
	}
	consumer, err := sarama.NewConsumer([]string{Config.Kafka.BootstrapServers}, Config.GetSaramaConfig())
	if err != nil {
		log.Fatal(err)
	}

	return &EventBus{
		eventRegistry: make(map[string]EventRegistration),
		SyncProducer:  syncProducer,
		Consumer:      consumer,
	}
}

func (eb EventBus) pushMessageToQueue(topic string, message []byte) error {
	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := eb.SyncProducer.SendMessage(producerMessage)
	if err != nil {
		return err
	}

	return nil
}

func (eb EventBus) StartConsumers(ctx context.Context) {
	var consumers []sarama.PartitionConsumer

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for topic, registration := range eb.eventRegistry {
		consumer, err := eb.Consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
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

	if err := eb.Consumer.Close(); err != nil {
		fmt.Printf("Error closing worker: %v\n", err)
	}
}
