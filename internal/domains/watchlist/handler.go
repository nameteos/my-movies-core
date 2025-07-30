package watchlist

import (
	"context"
	"event-driven-go/internal/shared"
	"fmt"
	"log"
)

type Handler struct{}

func (h *Handler) Handle(ctx context.Context, event shared.Event) error {
	switch e := event.(type) {
	case *MovieAddedToWatchlistEvent:
		return h.handleMovieAddedToWatchlist(ctx, e)
	default:
		return fmt.Errorf("unsupported event type: %T", event)
	}
}

func (h *Handler) CanHandle(eventType string) bool {
	return eventType == MovieAddedToWatchlistEventType
}

func (h *Handler) handleMovieAddedToWatchlist(ctx context.Context, event *MovieAddedToWatchlistEvent) error {
	log.Printf("🎬 WATCHLIST: User %s added '%s' (%d, %s) to watchlist",
		event.UserID, event.Title, event.Year, event.Genre)

	// todo handle with service

	log.Printf("✅ WATCHLIST: Successfully processed event %s", event.GetID())
	return nil
}
