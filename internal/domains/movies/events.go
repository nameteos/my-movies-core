package movies

import (
	"event-driven-go/internal/shared"
)

const (
	MovieCreatedEventType = "movies_movie_created"
	MovieUpdatedEventType = "movies_movie_updated"
	MovieDeletedEventType = "movies_movie_deleted"
)

type MovieCreatedEvent struct {
	shared.BaseEvent
	MovieID     string   `json:"movie_id"`
	Title       string   `json:"title"`
	Year        int      `json:"year"`
	Genre       []string `json:"genre"`
	Director    []string `json:"director"`
	Description string   `json:"description"`
}

func init() {
	shared.GlobalEventBus.RegisterEventType(MovieCreatedEventType, &MovieCreatedEvent{}, &Handler{})
	shared.GlobalEventBus.RegisterEventType(MovieUpdatedEventType, &MovieUpdatedEvent{}, &Handler{})
	shared.GlobalEventBus.RegisterEventType(MovieDeletedEventType, &MovieDeletedEvent{}, &Handler{})
}

func (e MovieCreatedEvent) GetPayload() interface{} {
	return struct {
		MovieID     string   `json:"movie_id"`
		Title       string   `json:"title"`
		Year        int      `json:"year"`
		Genre       []string `json:"genre"`
		Director    []string `json:"director"`
		Description string   `json:"description"`
	}{
		MovieID:     e.MovieID,
		Title:       e.Title,
		Year:        e.Year,
		Genre:       e.Genre,
		Director:    e.Director,
		Description: e.Description,
	}
}

func NewMovieCreatedEvent(movieID, title, description string, year int, genre, director []string) *MovieCreatedEvent {
	return &MovieCreatedEvent{
		BaseEvent:   shared.NewBaseEvent(MovieCreatedEventType),
		MovieID:     movieID,
		Title:       title,
		Year:        year,
		Genre:       genre,
		Director:    director,
		Description: description,
	}
}

type MovieUpdatedEvent struct {
	shared.BaseEvent
	MovieID string `json:"movie_id"`
	Title   string `json:"title"`
}

func (e MovieUpdatedEvent) GetPayload() interface{} {
	return struct {
		MovieID string `json:"movie_id"`
		Title   string `json:"title"`
	}{
		MovieID: e.MovieID,
		Title:   e.Title,
	}
}

func NewMovieUpdatedEvent(movieID, title string) *MovieUpdatedEvent {
	return &MovieUpdatedEvent{
		BaseEvent: shared.NewBaseEvent(MovieUpdatedEventType),
		MovieID:   movieID,
		Title:     title,
	}
}

type MovieDeletedEvent struct {
	shared.BaseEvent
	MovieID string `json:"movie_id"`
	Title   string `json:"title"`
}

func (e MovieDeletedEvent) GetPayload() interface{} {
	return struct {
		MovieID string `json:"movie_id"`
		Title   string `json:"title"`
	}{
		MovieID: e.MovieID,
		Title:   e.Title,
	}
}

func NewMovieDeletedEvent(movieID, title string) *MovieDeletedEvent {
	return &MovieDeletedEvent{
		BaseEvent: shared.NewBaseEvent(MovieDeletedEventType),
		MovieID:   movieID,
		Title:     title,
	}
}
