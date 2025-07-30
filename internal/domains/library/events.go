package library

import (
	"time"

	"event-driven-go/internal/shared"
)

// Event type constants
const (
	MovieWatchedEventType = "library_movie_watched"
)

type MovieWatchedEvent struct {
	shared.BaseEvent
	UserID    string    `json:"user_id"`
	MovieID   string    `json:"movie_id"`
	Title     string    `json:"title"`
	WatchedAt time.Time `json:"watched_at"`
	Duration  int       `json:"duration_minutes,omitempty"`
}

func init() {
	shared.GlobalEventBus.RegisterEventType(MovieWatchedEventType, &MovieWatchedEvent{}, &Handler{})
}

func (e MovieWatchedEvent) GetPayload() interface{} {
	return struct {
		UserID    string    `json:"user_id"`
		MovieID   string    `json:"movie_id"`
		Title     string    `json:"title"`
		WatchedAt time.Time `json:"watched_at"`
		Duration  int       `json:"duration_minutes,omitempty"`
	}{
		UserID:    e.UserID,
		MovieID:   e.MovieID,
		Title:     e.Title,
		WatchedAt: e.WatchedAt,
		Duration:  e.Duration,
	}
}

func NewMovieWatchedEvent(userID, movieID, title string, watchedAt time.Time, duration int) *MovieWatchedEvent {
	return &MovieWatchedEvent{
		BaseEvent: shared.NewBaseEvent(MovieWatchedEventType),
		UserID:    userID,
		MovieID:   movieID,
		Title:     title,
		WatchedAt: watchedAt,
		Duration:  duration,
	}
}
