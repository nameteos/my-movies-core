package movies

import (
	"context"
	"fmt"
	schema "github.com/nameteos/my-movies-db-schema/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	CreateMovie(ctx context.Context, movie *schema.Movie) (*schema.Movie, error)
	GetMovieByID(ctx context.Context, id string) (*schema.Movie, error)
	UpdateMovie(ctx context.Context, movie *schema.Movie) (*schema.Movie, error)
	DeleteMovie(ctx context.Context, id string) error
	SearchMovies(ctx context.Context, query string, limit, offset int) ([]*schema.Movie, error)
	GetMoviesByGenre(ctx context.Context, genre string, limit, offset int) ([]*schema.Movie, error)
	GetMoviesByYear(ctx context.Context, year int, limit, offset int) ([]*schema.Movie, error)
	GetMoviesByDirector(ctx context.Context, director string, limit, offset int) ([]*schema.Movie, error)
	GetRecentMovies(ctx context.Context, limit, offset int) ([]*schema.Movie, error)
}

type MongoRepository struct {
	collection *mongo.Collection
	indexer    *schema.MongoIndexer
}

func NewMongoIndexer(db *mongo.Database) *schema.MongoIndexer {
	return &schema.MongoIndexer{
		Collection:       db.Collection("movies"),
		VectorDimensions: 1536,
	}
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("movies"),
		indexer:    NewMongoIndexer(db),
	}
}

func (r *MongoRepository) CreateMovie(ctx context.Context, movie *schema.Movie) (*schema.Movie, error) {
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

func (r *MongoRepository) GetMovieByID(ctx context.Context, id string) (*schema.Movie, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	var movie schema.Movie
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}

	return &movie, nil
}

func (r *MongoRepository) UpdateMovie(ctx context.Context, movie *schema.Movie) (*schema.Movie, error) {
	movie.UpdatedAt = time.Now()

	filter := bson.M{"_id": movie.ID}
	update := bson.M{"$set": movie}

	result := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedMovie schema.Movie
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

func (r *MongoRepository) SearchMovies(ctx context.Context, query string, limit, offset int) ([]*schema.Movie, error) {
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

	var movies []*schema.Movie
	for cursor.Next(ctx) {
		var movie schema.Movie
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
func (r *MongoRepository) GetMoviesByGenre(ctx context.Context, genre string, limit, offset int) ([]*schema.Movie, error) {
	filter := bson.M{"genre": genre}
	return r.findMoviesWithFilter(ctx, filter, limit, offset)
}

func (r *MongoRepository) GetMoviesByYear(ctx context.Context, year int, limit, offset int) ([]*schema.Movie, error) {
	filter := bson.M{"year": year}
	return r.findMoviesWithFilter(ctx, filter, limit, offset)
}

func (r *MongoRepository) GetMoviesByDirector(ctx context.Context, director string, limit, offset int) ([]*schema.Movie, error) {
	filter := bson.M{"director": director}
	return r.findMoviesWithFilter(ctx, filter, limit, offset)
}

func (r *MongoRepository) GetRecentMovies(ctx context.Context, limit, offset int) ([]*schema.Movie, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent movies: %w", err)
	}
	defer cursor.Close(ctx)

	var movies []*schema.Movie
	for cursor.Next(ctx) {
		var movie schema.Movie
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

func (r *MongoRepository) findMoviesWithFilter(ctx context.Context, filter bson.M, limit, offset int) ([]*schema.Movie, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"year": -1, "title": 1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find movies: %w", err)
	}
	defer cursor.Close(ctx)

	var movies []*schema.Movie
	for cursor.Next(ctx) {
		var movie schema.Movie
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
	if _, err := r.indexer.CreateIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}
	return nil
}
