package user

import (
	"context"
	"fmt"
	"log"

	"event-driven-go/internal/shared"
)

type Handler struct{}

func (h *Handler) Handle(ctx context.Context, event shared.Event) error {
	switch e := event.(type) {
	case *UserRegisteredEvent:
		return h.handleUserRegistered(ctx, e)
	case *UserUpdatedEvent:
		return h.handleUserUpdated(ctx, e)
	case *UserDeletedEvent:
		return h.handleUserDeleted(ctx, e)
	default:
		return fmt.Errorf("unsupported event type: %T", event)
	}
}

func (h *Handler) CanHandle(eventType string) bool {
	return eventType == UserRegisteredEventType ||
		eventType == UserUpdatedEventType ||
		eventType == UserDeletedEventType
}

func (h *Handler) handleUserRegistered(ctx context.Context, event *UserRegisteredEvent) error {
	log.Printf("👤 USER: New user registered - '%s' (%s)",
		event.Username, event.Email)

	// todo service inject

	return nil
}

// handleUserUpdated processes UserUpdatedEvent
func (h *Handler) handleUserUpdated(ctx context.Context, event *UserUpdatedEvent) error {
	log.Printf("👤 USER: User '%s' updated their profile", event.Username)

	return nil
}

// handleUserDeleted processes UserDeletedEvent
func (h *Handler) handleUserDeleted(ctx context.Context, event *UserDeletedEvent) error {
	log.Printf("👤 USER: User '%s' has been deleted", event.Username)

	return nil
}
