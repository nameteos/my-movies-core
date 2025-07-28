package rating

import (
	"event-driven-go/internal/shared"
)

const (
	MovieRatedEventType   = "rating.movie_rated"
	MovieUnratedEventType = "rating.movie_unrated"
)

type MovieRatedEvent struct {
	shared.BaseEvent
	UserID  string  `json:"user_id"`
	MovieID string  `json:"movie_id"`
	Title   string  `json:"title"`
	Rating  float64 `json:"rating"`
	Review  string  `json:"review,omitempty"`
}

func (e MovieRatedEvent) GetPayload() interface{} {
	return struct {
		UserID  string  `json:"user_id"`
		MovieID string  `json:"movie_id"`
		Title   string  `json:"title"`
		Rating  float64 `json:"rating"`
		Review  string  `json:"review,omitempty"`
	}{
		UserID:  e.UserID,
		MovieID: e.MovieID,
		Title:   e.Title,
		Rating:  e.Rating,
		Review:  e.Review,
	}
}

func NewMovieRatedEvent(userID, movieID, title string, rating float64, review string) *MovieRatedEvent {
	return &MovieRatedEvent{
		BaseEvent: shared.NewBaseEvent(MovieRatedEventType),
		UserID:    userID,
		MovieID:   movieID,
		Title:     title,
		Rating:    rating,
		Review:    review,
	}
}

type MovieUnratedEvent struct {
	shared.BaseEvent
	UserID  string `json:"user_id"`
	MovieID string `json:"movie_id"`
	Title   string `json:"title"`
}

func (e MovieUnratedEvent) GetPayload() interface{} {
	return struct {
		UserID  string `json:"user_id"`
		MovieID string `json:"movie_id"`
		Title   string `json:"title"`
	}{
		UserID:  e.UserID,
		MovieID: e.MovieID,
		Title:   e.Title,
	}
}

func NewMovieUnratedEvent(userID, movieID, title string) *MovieUnratedEvent {
	return &MovieUnratedEvent{
		BaseEvent: shared.NewBaseEvent(MovieUnratedEventType),
		UserID:    userID,
		MovieID:   movieID,
		Title:     title,
	}
}
