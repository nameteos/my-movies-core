# my-movies-go

Learning project. POC for event driven architecture in Go with Kafka.

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd my-movies-go
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application locally with hot reloading:
```bash
cp docker-compose.override.yml.dist docker-compose.override.yml
docker-compose up -d
```

## Event examples

#### `watchlist.movie_added`
Triggered when a user adds a movie to their watchlist.
```json
{
  "id": "uuid",
  "type": "watchlist.movie_added",
  "timestamp": "2024-01-01T12:00:00Z",
  "user_id": "user-123",
  "movie_id": "movie-456",
  "title": "The Matrix",
  "genre": "Sci-Fi",
  "year": 1999
}
```
#### `library.movie_watched`
Triggered when a user marks a movie as watched.
```json
{
  "id": "uuid",
  "type": "library.movie_watched",
  "timestamp": "2024-01-01T12:00:00Z",
  "user_id": "user-123",
  "movie_id": "movie-456",
  "title": "The Matrix",
  "watched_at": "2024-01-01T10:00:00Z",
  "duration_minutes": 136
}
```
#### `rating.movie_rated`
Triggered when a user rates a movie.
```json
{
  "id": "uuid",
  "type": "rating.movie_rated",
  "timestamp": "2024-01-01T12:00:00Z",
  "user_id": "user-123",
  "movie_id": "movie-456",
  "title": "The Matrix",
  "rating": 4.5,
  "review": "Amazing movie!"
}
```
#### `rating.movie_unrated`
Triggered when a user removes their rating.
```json
{
  "id": "uuid",
  "type": "rating.movie_unrated",
  "timestamp": "2024-01-01T12:00:00Z",
  "user_id": "user-123",
  "movie_id": "movie-456",
  "title": "The Matrix"
}
```

