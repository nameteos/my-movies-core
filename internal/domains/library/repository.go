package library

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddWatchHistory(ctx context.Context, userID, movieID string, watchedAt time.Time, duration int) (*WatchHistory, error) {
	history := &WatchHistory{
		ID:        uuid.New().String(),
		UserID:    userID,
		MovieID:   movieID,
		WatchedAt: watchedAt,
		Duration:  duration,
	}

	result := r.db.WithContext(ctx).Create(history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to add watch history: %w", result.Error)
	}

	return history, nil
}

func (r *Repository) GetUserWatchHistory(ctx context.Context, userID string, limit, offset int) ([]*WatchHistory, error) {
	var history []*WatchHistory

	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("watched_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&history)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get watch history: %w", result.Error)
	}

	return history, nil
}

func (r *Repository) GetMovieWatchCount(ctx context.Context, movieID string) (int, error) {
	var count int64

	result := r.db.WithContext(ctx).
		Model(&WatchHistory{}).
		Where("movie_id = ?", movieID).
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to get movie watch count: %w", result.Error)
	}

	return int(count), nil
}

func (r *Repository) HasUserWatchedMovie(ctx context.Context, userID, movieID string) (bool, error) {
	var count int64

	result := r.db.WithContext(ctx).
		Model(&WatchHistory{}).
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check watch history: %w", result.Error)
	}

	return count > 0, nil
}

func (r *Repository) GetRecentlyWatchedMovies(ctx context.Context, limit int) ([]*WatchHistory, error) {
	var history []*WatchHistory

	// this is across all users
	result := r.db.WithContext(ctx).
		Order("watched_at DESC").
		Limit(limit).
		Find(&history)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get recently watched movies: %w", result.Error)
	}

	return history, nil
}

func (r *Repository) GetWatchingStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalMoviesWatched int64
	if err := r.db.WithContext(ctx).
		Model(&WatchHistory{}).
		Where("user_id = ?", userID).
		Count(&totalMoviesWatched).Error; err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}
	stats["total_movies_watched"] = totalMoviesWatched

	var totalWatchTime int64
	if err := r.db.WithContext(ctx).
		Model(&WatchHistory{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(duration), 0)").
		Scan(&totalWatchTime).Error; err != nil {
		return nil, fmt.Errorf("failed to get total duration: %w", err)
	}
	stats["total_minutes"] = totalWatchTime
	stats["total_hours"] = float64(totalWatchTime) / 60.0

	var totalMoviesWatchedThisMonth int64
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	if err := r.db.WithContext(ctx).
		Model(&WatchHistory{}).
		Where("user_id = ? AND watched_at >= ?", userID, startOfMonth).
		Count(&totalMoviesWatchedThisMonth).Error; err != nil {
		return nil, fmt.Errorf("failed to get this month count: %w", err)
	}
	stats["total_movies_watched_this_month"] = totalMoviesWatchedThisMonth

	return stats, nil
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&WatchHistory{})
}
