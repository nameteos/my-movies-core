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
	MovieID string `json:"movie_id"`
	Title   string `json:"title"`
}

func init() {
	shared.GlobalEventBus.RegisterEventType(MovieCreatedEventType, &MovieCreatedEvent{}, &Handler{})
	shared.GlobalEventBus.RegisterEventType(MovieUpdatedEventType, &MovieUpdatedEvent{}, &Handler{})
	shared.GlobalEventBus.RegisterEventType(MovieDeletedEventType, &MovieDeletedEvent{}, &Handler{})
}

func (e MovieCreatedEvent) GetPayload() interface{} {
	return struct {
		MovieID string `json:"movie_id"`
		Title   string `json:"title"`
	}{
		MovieID: e.MovieID,
		Title:   e.Title,
	}
}

func NewMovieCreatedEvent(movieID string, title string) *MovieCreatedEvent {
	return &MovieCreatedEvent{
		BaseEvent: shared.NewBaseEvent(MovieCreatedEventType),
		MovieID:   movieID,
		Title:     title,
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
