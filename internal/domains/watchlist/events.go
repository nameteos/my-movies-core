package watchlist

import "event-driven-go/internal/shared"

const (
	MovieAddedToWatchlistEventType = "watchlist_movie_added"
)

func init() {
	shared.GlobalEventBus.RegisterEventType(MovieAddedToWatchlistEventType, &MovieAddedToWatchlistEvent{}, &Handler{})
}

type MovieAddedToWatchlistEvent struct {
	shared.BaseEvent
	UserID  string `json:"user_id"`
	MovieID string `json:"movie_id"`
	Title   string `json:"title"`
}

func (e MovieAddedToWatchlistEvent) GetPayload() interface{} {
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

func NewMovieAddedToWatchlistEvent(userID string, movieID string, title string) *MovieAddedToWatchlistEvent {
	return &MovieAddedToWatchlistEvent{
		BaseEvent: shared.NewBaseEvent(MovieAddedToWatchlistEventType),
		UserID:    userID,
		MovieID:   movieID,
		Title:     title,
	}
}
