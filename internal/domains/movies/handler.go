package movies

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
	return &Handler{logger}
}

func (h *Handler) Handle(ctx context.Context, event shared.Event) error {
	switch e := event.(type) {
	case *MovieCreatedEvent:
		return h.handleMovieCreated(ctx, e)
	case *MovieUpdatedEvent:
		return h.handleMovieUpdated(ctx, e)
	case *MovieDeletedEvent:
		return h.handleMovieDeleted(ctx, e)
	default:
		return fmt.Errorf("unsupported event type: %T", event)
	}
}

func (h *Handler) CanHandle(eventType string) bool {
	return eventType == MovieCreatedEventType ||
		eventType == MovieUpdatedEventType ||
		eventType == MovieDeletedEventType
}

func (h *Handler) handleMovieCreated(ctx context.Context, event *MovieCreatedEvent) error {
	h.logger.Printf("🎬 MOVIES: New movie '%s' (%d) added to catalog",
		event.Title, event.Year)

	// todo index movie

	return nil
}

// handleMovieUpdated processes MovieUpdatedEvent
func (h *Handler) handleMovieUpdated(ctx context.Context, event *MovieUpdatedEvent) error {
	h.logger.Printf("🎬 MOVIES: Movie '%s' has been updated", event.Title)

	return nil
}

// handleMovieDeleted processes MovieDeletedEvent
func (h *Handler) handleMovieDeleted(ctx context.Context, event *MovieDeletedEvent) error {
	h.logger.Printf("🎬 MOVIES: Movie '%s' has been deleted from catalog", event.Title)
	return nil
}
