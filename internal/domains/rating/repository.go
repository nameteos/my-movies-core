package rating

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddRating(ctx context.Context, userID, movieID string, rating float64, review string) (*MovieRating, error) {
	movieRating := &MovieRating{
		ID:      uuid.New().String(),
		UserID:  userID,
		MovieID: movieID,
		Rating:  rating,
		Review:  review,
	}

	result := r.db.WithContext(ctx).Create(movieRating)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to add rating: %w", result.Error)
	}

	return movieRating, nil
}

func (r *Repository) UpdateRating(ctx context.Context, userID, movieID string, rating float64, review string) (*MovieRating, error) {
	var movieRating MovieRating

	// Find existing rating
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		First(&movieRating)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("rating not found")
		}
		return nil, fmt.Errorf("failed to find rating: %w", result.Error)
	}

	// Update fields
	movieRating.Rating = rating
	movieRating.Review = review

	if err := r.db.WithContext(ctx).Save(&movieRating).Error; err != nil {
		return nil, fmt.Errorf("failed to update rating: %w", err)
	}

	return &movieRating, nil
}

func (r *Repository) UpsertRating(ctx context.Context, userID, movieID string, rating float64, review string) (*MovieRating, error) {
	var movieRating MovieRating

	// Try to find existing rating
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		First(&movieRating)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new rating
		return r.AddRating(ctx, userID, movieID, rating, review)
	} else if result.Error != nil {
		return nil, fmt.Errorf("failed to check existing rating: %w", result.Error)
	}

	// Update existing rating
	return r.UpdateRating(ctx, userID, movieID, rating, review)
}

func (r *Repository) RemoveRating(ctx context.Context, userID, movieID string) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		Delete(&MovieRating{})

	if result.Error != nil {
		return fmt.Errorf("failed to remove rating: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("rating not found")
	}

	return nil
}

func (r *Repository) GetUserRating(ctx context.Context, userID, movieID string) (*MovieRating, error) {
	var rating MovieRating

	result := r.db.WithContext(ctx).
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		First(&rating)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("rating not found")
		}
		return nil, fmt.Errorf("failed to get user rating: %w", result.Error)
	}

	return &rating, nil
}

func (r *Repository) GetMovieRatings(ctx context.Context, movieID string, limit, offset int) ([]*MovieRating, error) {
	var ratings []*MovieRating

	result := r.db.WithContext(ctx).
		Where("movie_id = ?", movieID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&ratings)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get movie ratings: %w", result.Error)
	}

	return ratings, nil
}

func (r *Repository) GetMovieAverageRating(ctx context.Context, movieID string) (float64, int, error) {
	var result struct {
		AvgRating float64
		Count     int64
	}

	err := r.db.WithContext(ctx).
		Model(&MovieRating{}).
		Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as count").
		Where("movie_id = ?", movieID).
		Scan(&result).Error

	if err != nil {
		return 0, 0, fmt.Errorf("failed to get movie average rating: %w", err)
	}

	return result.AvgRating, int(result.Count), nil
}

func (r *Repository) GetUserRatings(ctx context.Context, userID string, limit, offset int) ([]*MovieRating, error) {
	var ratings []*MovieRating

	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&ratings)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user ratings: %w", result.Error)
	}

	return ratings, nil
}

func (r *Repository) GetTopRatedMovies(ctx context.Context, limit int) ([]struct {
	MovieID   string  `json:"movie_id"`
	AvgRating float64 `json:"avg_rating"`
	Count     int64   `json:"count"`
}, error) {
	var results []struct {
		MovieID   string  `json:"movie_id"`
		AvgRating float64 `json:"avg_rating"`
		Count     int64   `json:"count"`
	}

	err := r.db.WithContext(ctx).
		Model(&MovieRating{}).
		Select("movie_id, AVG(rating) as avg_rating, COUNT(*) as count").
		Group("movie_id").
		Having("COUNT(*) >= ?", 3). // Only movies with at least 3 ratings
		Order("avg_rating DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top rated movies: %w", err)
	}

	return results, nil
}

func (r *Repository) GetRatingDistribution(ctx context.Context, movieID string) (map[string]int64, error) {
	var results []struct {
		Rating string `json:"rating"`
		Count  int64  `json:"count"`
	}

	err := r.db.WithContext(ctx).
		Model(&MovieRating{}).
		Select("FLOOR(rating) as rating, COUNT(*) as count").
		Where("movie_id = ?", movieID).
		Group("FLOOR(rating)").
		Order("rating").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get rating distribution: %w", err)
	}

	distribution := make(map[string]int64)
	for _, result := range results {
		distribution[result.Rating+" stars"] = result.Count
	}

	return distribution, nil
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&MovieRating{})
}
