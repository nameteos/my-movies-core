package handlers

import (
	"context"
	"log"

	"event-driven-go/internal/domains/library"
	"event-driven-go/internal/domains/rating"
	"event-driven-go/internal/domains/watchlist"
	"event-driven-go/internal/shared"
)

type NotificationHandler struct {
	logger *log.Logger
}

func NewNotificationHandler(logger *log.Logger) *NotificationHandler {
	if logger == nil {
		logger = log.Default()
	}
	return &NotificationHandler{logger: logger}
}

func (h *NotificationHandler) Handle(ctx context.Context, event shared.Event) error {
	switch e := event.(type) {
	case *watchlist.MovieAddedToWatchlistEvent:
		h.logger.Printf("📱 NOTIFICATION: Movie '%s' added to watchlist for user %s", e.Title, e.UserID)
	case *library.MovieWatchedEvent:
		h.logger.Printf("📱 NOTIFICATION: User %s watched '%s'", e.UserID, e.Title)
	case *rating.MovieRatedEvent:
		h.logger.Printf("📱 NOTIFICATION: User %s rated '%s' %.1f stars", e.UserID, e.Title, e.Rating)
	case *rating.MovieUnratedEvent:
		h.logger.Printf("📱 NOTIFICATION: User %s removed rating for '%s'", e.UserID, e.Title)
	default:
		h.logger.Printf("📱 NOTIFICATION: Unknown event type %s", event.GetType())
	}

	// Here you would typically:
	// - Send push notifications to mobile apps
	// - Send email notifications if enabled
	// - Update social feeds
	// - Send notifications to followers
	// - Update notification center/inbox

	return nil
}

// CanHandle checks if this handler can process the given event type
func (h *NotificationHandler) CanHandle(eventType string) bool {
	return eventType == watchlist.MovieAddedToWatchlistEventType ||
		eventType == library.MovieWatchedEventType ||
		eventType == rating.MovieRatedEventType ||
		eventType == rating.MovieUnratedEventType
}
