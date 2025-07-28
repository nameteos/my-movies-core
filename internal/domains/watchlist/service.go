package watchlist

import (
	"context"
	"fmt"
	"strings"

	"event-driven-go/internal/domains/movies"
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

func (s *Service) AddMovie(ctx context.Context, userID string, movie *movies.Movie) error {
	// Validate input
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	if movie == nil {
		return fmt.Errorf("movie cannot be nil")
	}

	// todo logic

	movieID := movie.IDString()              // Convert ObjectID to string
	genre := strings.Join(movie.Genre, ", ") // Convert []string to single string for event

	event := NewMovieAddedToWatchlistEvent(
		userID,
		movieID,
		movie.Title,
		genre,
		movie.Year,
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
