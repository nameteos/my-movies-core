package library

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
	case *MovieWatchedEvent:
		return h.handleMovieWatched(ctx, e)
	default:
		return fmt.Errorf("unsupported event type: %T", event)
	}
}

func (h *Handler) CanHandle(eventType string) bool {
	return eventType == MovieWatchedEventType
}

func (h *Handler) handleMovieWatched(ctx context.Context, event *MovieWatchedEvent) error {
	durationText := ""
	if event.Duration > 0 {
		durationText = fmt.Sprintf(" (Duration: %d min)", event.Duration)
	}

	h.logger.Printf("📚 LIBRARY: User %s watched '%s' at %s%s",
		event.UserID, event.Title, event.WatchedAt.Format("2006-01-02 15:04:05"), durationText)

	// todo store move in DB, do other stuff

	return nil
}
