package library

import (
	"context"
	"fmt"
	"time"

	"event-driven-go/internal/shared"
)

type Service struct {
	eventBus *shared.EventBus
}

func NewService(eventBus *shared.EventBus) *Service {
	return &Service{
		eventBus: eventBus,
	}
}

func (s *Service) MarkAsWatched(ctx context.Context, userID, movieID, title string, watchedAt time.Time, duration int) error {
	// todo implement validator?
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	if movieID == "" {
		return fmt.Errorf("movie ID cannot be empty")
	}
	if title == "" {
		return fmt.Errorf("movie title cannot be empty")
	}

	// todo:
	// - Store watch history in database
	// - listener?- Remove from watchlist if present
	// - listener?- Update user's viewing statistics

	// Create and publish the event
	event := NewMovieWatchedEvent(userID, movieID, title, watchedAt, duration)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		return fmt.Errorf("failed to publish movie watched event: %w", err)
	}

	return nil
}

func (s *Service) GetWatchHistory(ctx context.Context, userID string) error {

	// TODO: Implement
	return nil
}
