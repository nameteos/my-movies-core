package watchlist

import (
	"context"
	"fmt"

	"event-driven-go/internal/shared"
	schema "github.com/nameteos/my-movies-db-schema/mongodb"
)

type Service struct {
	eventBus *shared.EventBus
}

func NewService(eventBus *shared.EventBus) *Service {
	return &Service{
		eventBus: eventBus,
	}
}

func (s *Service) AddMovie(ctx context.Context, userID string, movie *schema.Movie) error {
	// Validate input
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	if movie == nil {
		return fmt.Errorf("movie cannot be nil")
	}

	// todo logic

	movieID := movie.ID.Hex() // Convert ObjectID to string

	event := NewMovieAddedToWatchlistEvent(
		userID,
		movieID,
		movie.Title,
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		return fmt.Errorf("failed to publish watchlist event: %w", err)
	}

	return nil
}

// RemoveMovie removes a movie from a user's watchlist
func (s *Service) RemoveMovie(ctx context.Context, userID, movieID string) error {
	// Validate input
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	if movieID == "" {
		return fmt.Errorf("movie ID cannot be empty")
	}

	// TODO: Implement removal logic and event
	return nil
}
