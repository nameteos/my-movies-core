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

3. Run the application locally:
```bash
go run cmd/main.go
```

## Docker Development Setup

This project includes a development-friendly Docker setup that allows you to make code changes without rebuilding the Docker image.

### Features

- **Live Reloading**: The application automatically rebuilds and restarts when you make code changes
- **Volume Mounting**: Your local codebase is mounted into the container, so changes are immediately available
- **Development Tools**: The development container includes tools for debugging and development

### Running with Docker

1. Start the development environment:
```bash
docker-compose -f docker-compose.yml -f docker-compose.override.yml up
```

2. Make changes to your code locally, and they will be automatically detected and the application will rebuild

3. To stop the development environment:
```bash
docker-compose -f docker-compose.yml -f docker-compose.override.yml down
```

### How It Works

The development setup uses:
- A `docker-compose.override.yml` file that extends the main docker-compose.yml
- A `Dockerfile.dev` optimized for development
- The [Air](https://github.com/cosmtrek/air) tool for live reloading

### Troubleshooting

- **Changes not being detected**: Make sure you're using both docker-compose.yml and docker-compose.override.yml files
- **Build errors**: Check the container logs for compilation errors
- **Connection issues**: The development setup preserves all the connection settings from the main docker-compose.yml
- **"air: executable not found"**: This issue has been fixed by using the full path to the air executable (/go/bin/air) in the docker-compose.override.yml file

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

