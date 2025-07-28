package rating

import (
	"context"
	"fmt"

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

func (s *Service) RateMovie(ctx context.Context, userID, movieID, title string, rating float64, review string) error {
	// todo use validator for all service methods?
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	if movieID == "" {
		return fmt.Errorf("movie ID cannot be empty")
	}
	if title == "" {
		return fmt.Errorf("movie title cannot be empty")
	}
	if rating < 0 || rating > 5 {
		return fmt.Errorf("rating must be between 0 and 5, got %.1f", rating)
	}

	// Create and publish the event
	event := NewMovieRatedEvent(userID, movieID, title, rating, review)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		return fmt.Errorf("failed to publish movie rating event: %w", err)
	}

	return nil
}

func (s *Service) RemoveRating(ctx context.Context, userID, movieID, title string) error {
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	if movieID == "" {
		return fmt.Errorf("movie ID cannot be empty")
	}
	if title == "" {
		return fmt.Errorf("movie title cannot be empty")
	}

	event := NewMovieUnratedEvent(userID, movieID, title)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		return fmt.Errorf("failed to publish movie unrated event: %w", err)
	}

	return nil
}

func (s *Service) GetMovieRatings(ctx context.Context, movieID string) error {
	return nil
}
