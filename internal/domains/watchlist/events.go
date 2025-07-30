package watchlist

import "event-driven-go/internal/shared"

const (
	MovieAddedToWatchlistEventType = "watchlist.movie_added"
)

func init() {
	shared.GlobalEventBus.RegisterEventType(MovieAddedToWatchlistEventType, &MovieAddedToWatchlistEvent{}, &Handler{})
}

type MovieAddedToWatchlistEvent struct {
	shared.BaseEvent
	UserID  string `json:"user_id"`
	MovieID string `json:"movie_id"`
	Title   string `json:"title"`
	Genre   string `json:"genre"`
	Year    int    `json:"year"`
}

func (e MovieAddedToWatchlistEvent) GetPayload() interface{} {
	return struct {
		UserID  string `json:"user_id"`
		MovieID string `json:"movie_id"`
		Title   string `json:"title"`
		Genre   string `json:"genre"`
		Year    int    `json:"year"`
	}{
		UserID:  e.UserID,
		MovieID: e.MovieID,
		Title:   e.Title,
		Genre:   e.Genre,
		Year:    e.Year,
	}
}

func NewMovieAddedToWatchlistEvent(userID, movieID, title, genre string, year int) *MovieAddedToWatchlistEvent {
	return &MovieAddedToWatchlistEvent{
		BaseEvent: shared.NewBaseEvent(MovieAddedToWatchlistEventType),
		UserID:    userID,
		MovieID:   movieID,
		Title:     title,
		Genre:     genre,
		Year:      year,
	}
}
