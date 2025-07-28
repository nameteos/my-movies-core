package movies

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	CreateMovie(ctx context.Context, movie *Movie) (*Movie, error)
	GetMovieByID(ctx context.Context, id string) (*Movie, error)
	UpdateMovie(ctx context.Context, movie *Movie) (*Movie, error)
	DeleteMovie(ctx context.Context, id string) error
	SearchMovies(ctx context.Context, query string, limit, offset int) ([]*Movie, error)
	GetMoviesByGenre(ctx context.Context, genre string, limit, offset int) ([]*Movie, error)
	GetMoviesByYear(ctx context.Context, year int, limit, offset int) ([]*Movie, error)
	GetMoviesByDirector(ctx context.Context, director string, limit, offset int) ([]*Movie, error)
	GetRecentMovies(ctx context.Context, limit, offset int) ([]*Movie, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("movies"),
	}
}

func (r *MongoRepository) CreateMovie(ctx context.Context, movie *Movie) (*Movie, error) {
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	if movie.ID.IsZero() {
		movie.ID = primitive.NewObjectID()
	}

	_, err := r.collection.InsertOne(ctx, movie)
	if err != nil {
		return nil, fmt.Errorf("failed to create movie: %w", err)
	}

	return movie, nil
}

func (r *MongoRepository) GetMovieByID(ctx context.Context, id string) (*Movie, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	var movie Movie
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}

	return &movie, nil
}

func (r *MongoRepository) UpdateMovie(ctx context.Context, movie *Movie) (*Movie, error) {
	movie.UpdatedAt = time.Now()

	filter := bson.M{"_id": movie.ID}
	update := bson.M{"$set": movie}

	result := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedMovie Movie
	if err := result.Decode(&updatedMovie); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, fmt.Errorf("failed to update movie: %w", err)
	}

	return &updatedMovie, nil
}

func (r *MongoRepository) DeleteMovie(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid movie ID: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("movie not found")
	}

	return nil
}

func (r *MongoRepository) SearchMovies(ctx context.Context, query string, limit, offset int) ([]*Movie, error) {
	// Create text search filter
	filter := bson.M{
		"$text": bson.M{
			"$search": query,
		},
	}

	// Set up options for pagination and sorting by text score
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"score": bson.M{"$meta": "textScore"}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer cursor.Close(ctx)

	var movies []*Movie
	for cursor.Next(ctx) {
		var movie Movie
		if err := cursor.Decode(&movie); err != nil {
			return nil, fmt.Errorf("failed to decode movie: %w", err)
		}
		movies = append(movies, &movie)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return movies, nil
}

// todo use cursor
func (r *MongoRepository) GetMoviesByGenre(ctx context.Context, genre string, limit, offset int) ([]*Movie, error) {
	filter := bson.M{"genre": genre}
	return r.findMoviesWithFilter(ctx, filter, limit, offset)
}

func (r *MongoRepository) GetMoviesByYear(ctx context.Context, year int, limit, offset int) ([]*Movie, error) {
	filter := bson.M{"year": year}
	return r.findMoviesWithFilter(ctx, filter, limit, offset)
}

func (r *MongoRepository) GetMoviesByDirector(ctx context.Context, director string, limit, offset int) ([]*Movie, error) {
	filter := bson.M{"director": director}
	return r.findMoviesWithFilter(ctx, filter, limit, offset)
}

func (r *MongoRepository) GetRecentMovies(ctx context.Context, limit, offset int) ([]*Movie, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent movies: %w", err)
	}
	defer cursor.Close(ctx)

	var movies []*Movie
	for cursor.Next(ctx) {
		var movie Movie
		if err := cursor.Decode(&movie); err != nil {
			return nil, fmt.Errorf("failed to decode movie: %w", err)
		}
		movies = append(movies, &movie)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return movies, nil
}

func (r *MongoRepository) findMoviesWithFilter(ctx context.Context, filter bson.M, limit, offset int) ([]*Movie, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"year": -1, "title": 1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find movies: %w", err)
	}
	defer cursor.Close(ctx)

	var movies []*Movie
	for cursor.Next(ctx) {
		var movie Movie
		if err := cursor.Decode(&movie); err != nil {
			return nil, fmt.Errorf("failed to decode movie: %w", err)
		}
		movies = append(movies, &movie)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return movies, nil
}

func (r *MongoRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		// Text index for search
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "description", Value: "text"},
				{Key: "cast.name", Value: "text"},
				{Key: "crew.name", Value: "text"},
			},
		},
		// Compound index for genre and year
		{
			Keys: bson.D{
				{Key: "genre", Value: 1},
				{Key: "year", Value: -1},
			},
		},
		// Index for director queries
		{
			Keys: bson.D{
				{Key: "director", Value: 1},
			},
		},
		// Index for year queries
		{
			Keys: bson.D{
				{Key: "year", Value: -1},
			},
		},
		// Index for recently added movies
		{
			Keys: bson.D{
				{Key: "created_at", Value: -1},
			},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}
