package shared

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Event interface {
	GetID() string
	GetType() string
	GetTimestamp() time.Time
	GetPayload() interface{}
}

// BaseEvent provides common fields for all events
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

// EventBus manages event publishing and subscription
type EventBus struct {
	handlers map[string][]EventHandler
	mutex    sync.RWMutex
	logger   *log.Logger
}

// NewEventBus creates a new event bus
func NewEventBus(logger *log.Logger) *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
		logger:   logger,
	}
}

func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	eb.logger.Printf("Subscribed handler for event type: %s", eventType)
}

// Publish sends an event to all registered handlers
func (eb *EventBus) Publish(ctx context.Context, event Event) error {
	eb.mutex.RLock()
	handlers, exists := eb.handlers[event.GetType()]
	eb.mutex.RUnlock()

	if !exists {
		eb.logger.Printf("No handlers registered for event type: %s", event.GetType())
		return nil
	}

	eb.logger.Printf("Publishing event: %s (ID: %s)", event.GetType(), event.GetID())

	for _, handler := range handlers {
		if handler.CanHandle(event.GetType()) {
			if err := handler.Handle(ctx, event); err != nil {
				eb.logger.Printf("Error handling event %s: %v", event.GetID(), err)
				return fmt.Errorf("handler error for event %s: %w", event.GetID(), err)
			}
		}
	}

	return nil
}

// NewBaseEvent creates a new base event with generated ID and current timestamp
func NewBaseEvent(eventType string) BaseEvent {
	return BaseEvent{
		ID:        uuid.New().String(),
		Type:      eventType,
		Timestamp: time.Now(),
	}
}
