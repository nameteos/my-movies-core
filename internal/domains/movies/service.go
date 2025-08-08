package movies

import (
	"context"
	"fmt"

	"event-driven-go/internal/shared"
	schema "github.com/nameteos/my-movies-db-schema/mongodb"
)

type Service struct {
	repository Repository
	eventBus   *shared.EventBus
}

func NewService(repository Repository, eventBus *shared.EventBus) *Service {
	return &Service{
		repository: repository,
		eventBus:   eventBus,
	}
}

func (s *Service) CreateMovie(ctx context.Context, movie *schema.Movie) (*schema.Movie, error) {
	// Validate input
	if movie == nil {
		return nil, fmt.Errorf("movie cannot be nil")
	}
	if movie.Title == "" {
		return nil, fmt.Errorf("movie title cannot be empty")
	}

	createdMovie, err := s.repository.CreateMovie(ctx, movie)
	if err != nil {
		return nil, fmt.Errorf("failed to create movie: %w", err)
	}

	event := NewMovieCreatedEvent(
		createdMovie.ID.Hex(),
		createdMovie.Title,
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation since movie was created
		// In a real system, you might want to implement compensation
		fmt.Printf("Warning: failed to publish movie created event: %v\n", err)
	}

	return createdMovie, nil
}

// GetMovieByID retrieves a movie by its ID
func (s *Service) GetMovieByID(ctx context.Context, id string) (*schema.Movie, error) {
	if id == "" {
		return nil, fmt.Errorf("movie ID cannot be empty")
	}

	return s.repository.GetMovieByID(ctx, id)
}

// UpdateMovie updates an existing movie and publishes an event
func (s *Service) UpdateMovie(ctx context.Context, movie *schema.Movie) (*schema.Movie, error) {
	// Validate input
	if movie == nil {
		return nil, fmt.Errorf("movie cannot be nil")
	}

	updatedMovie, err := s.repository.UpdateMovie(ctx, movie)
	if err != nil {
		return nil, fmt.Errorf("failed to update movie: %w", err)
	}

	// Publish movie updated event
	event := NewMovieUpdatedEvent(updatedMovie.ID.Hex(), updatedMovie.Title)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to publish movie updated event: %v\n", err)
	}

	return updatedMovie, nil
}

// DeleteMovie deletes a movie and publishes an event
func (s *Service) DeleteMovie(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("movie ID cannot be empty")
	}

	// Get movie first to get title for event
	movie, err := s.repository.GetMovieByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get movie before deletion: %w", err)
	}

	// Delete movie from repository
	if err := s.repository.DeleteMovie(ctx, id); err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	// Publish movie deleted event
	event := NewMovieDeletedEvent(id, movie.Title)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation since movie was deleted
		fmt.Printf("Warning: failed to publish movie deleted event: %v\n", err)
	}

	return nil
}

// SearchMovies searches for movies using various criteria
func (s *Service) SearchMovies(ctx context.Context, query string, limit, offset int) ([]*schema.Movie, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.repository.SearchMovies(ctx, query, limit, offset)
}

// GetMoviesByGenre retrieves movies by genre
func (s *Service) GetMoviesByGenre(ctx context.Context, genre string, limit, offset int) ([]*schema.Movie, error) {
	if genre == "" {
		return nil, fmt.Errorf("genre cannot be empty")
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.repository.GetMoviesByGenre(ctx, genre, limit, offset)
}

// GetMoviesByYear retrieves movies by year
func (s *Service) GetMoviesByYear(ctx context.Context, year int, limit, offset int) ([]*schema.Movie, error) {
	if year <= 0 {
		return nil, fmt.Errorf("year must be valid")
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.repository.GetMoviesByYear(ctx, year, limit, offset)
}

// GetMoviesByDirector retrieves movies by director
func (s *Service) GetMoviesByDirector(ctx context.Context, director string, limit, offset int) ([]*schema.Movie, error) {
	if director == "" {
		return nil, fmt.Errorf("director cannot be empty")
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.repository.GetMoviesByDirector(ctx, director, limit, offset)
}

// GetRecentMovies retrieves recently added movies
func (s *Service) GetRecentMovies(ctx context.Context, limit, offset int) ([]*schema.Movie, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.repository.GetRecentMovies(ctx, limit, offset)
}
