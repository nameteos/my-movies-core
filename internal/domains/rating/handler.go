package rating

import (
	"context"
	"fmt"
	"log"

	"event-driven-go/internal/shared"
)

type Handler struct{}

func (h *Handler) Handle(ctx context.Context, event shared.Event) error {
	switch e := event.(type) {
	case *MovieRatedEvent:
		return h.handleMovieRated(ctx, e)
	case *MovieUnratedEvent:
		return h.handleMovieUnrated(ctx, e)
	default:
		return fmt.Errorf("unsupported event type: %T", event)
	}
}

func (h *Handler) CanHandle(eventType string) bool {
	return eventType == MovieRatedEventType || eventType == MovieUnratedEventType
}

func (h *Handler) handleMovieRated(ctx context.Context, event *MovieRatedEvent) error {
	reviewText := ""
	if event.Review != "" {
		reviewText = fmt.Sprintf(" with review: \"%s\"", event.Review)
	}

	log.Printf("⭐ RATING: User %s rated '%s' %.1f/5%s",
		event.UserID, event.Title, event.Rating, reviewText)

	// todo service

	log.Printf("✅ RATING: Successfully processed rating event %s", event.GetID())
	return nil
}

// handleMovieUnrated processes MovieUnratedEvent
func (h *Handler) handleMovieUnrated(ctx context.Context, event *MovieUnratedEvent) error {
	log.Printf("❌ RATING: User %s removed rating for '%s'",
		event.UserID, event.Title)

	// todo service

	log.Printf("✅ RATING: Successfully processed unrating event %s", event.GetID())
	return nil
}
