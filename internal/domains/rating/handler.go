package rating

import (
	"context"
	"fmt"
	"log"

	"event-driven-go/internal/shared"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{logger: logger}
}

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

	h.logger.Printf("⭐ RATING: User %s rated '%s' %.1f/5%s",
		event.UserID, event.Title, event.Rating, reviewText)

	// todo service

	h.logger.Printf("✅ RATING: Successfully processed rating event %s", event.GetID())
	return nil
}

// handleMovieUnrated processes MovieUnratedEvent
func (h *Handler) handleMovieUnrated(ctx context.Context, event *MovieUnratedEvent) error {
	h.logger.Printf("❌ RATING: User %s removed rating for '%s'",
		event.UserID, event.Title)

	// todo service

	h.logger.Printf("✅ RATING: Successfully processed unrating event %s", event.GetID())
	return nil
}
