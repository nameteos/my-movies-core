package watchlist

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

func (r *Repository) AddToWatchlist(ctx context.Context, userID, movieID string, notes string) (*WatchlistEntry, error) {
	entry := &WatchlistEntry{
		ID:      uuid.New().String(),
		UserID:  userID,
		MovieID: movieID,
		AddedAt: time.Now(),
		Notes:   notes,
	}

	// GORM will handle the upsert with Clauses
	result := r.db.WithContext(ctx).
		Clauses().
		Create(entry)

	if result.Error != nil {
		// Handle unique constraint violation (user already has this movie in watchlist)
		if isDuplicateKeyError(result.Error) {
			// Update existing entry
			var existing WatchlistEntry
			if err := r.db.WithContext(ctx).
				Where("user_id = ? AND movie_id = ?", userID, movieID).
				First(&existing).Error; err != nil {
				return nil, fmt.Errorf("failed to find existing entry: %w", err)
			}

			existing.Notes = notes
			existing.AddedAt = time.Now()

			if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
				return nil, fmt.Errorf("failed to update existing entry: %w", err)
			}

			return &existing, nil
		}
		return nil, fmt.Errorf("failed to add to watchlist: %w", result.Error)
	}

	return entry, nil
}

func (r *Repository) RemoveFromWatchlist(ctx context.Context, userID, movieID string) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		Delete(&WatchlistEntry{})

	if result.Error != nil {
		return fmt.Errorf("failed to remove from watchlist: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("movie not found in watchlist")
	}

	return nil
}

func (r *Repository) GetUserWatchlist(ctx context.Context, userID string) ([]*WatchlistEntry, error) {
	var entries []*WatchlistEntry

	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("added_at DESC").
		Find(&entries)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get watchlist: %w", result.Error)
	}

	return entries, nil
}

func (r *Repository) IsInWatchlist(ctx context.Context, userID, movieID string) (bool, error) {
	var count int64

	result := r.db.WithContext(ctx).
		Model(&WatchlistEntry{}).
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check watchlist: %w", result.Error)
	}

	return count > 0, nil
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&WatchlistEntry{})
}

func isDuplicateKeyError(err error) bool {
	// PostgreSQL duplicate key error contains "duplicate key value violates unique constraint"
	// You might want to use a more robust error type checking here
	return err != nil && gorm.ErrDuplicatedKey == err
}
