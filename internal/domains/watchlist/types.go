package watchlist

import (
	"time"

	"gorm.io/gorm"
)

type WatchlistEntry struct {
	ID      string    `json:"id" db:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID  string    `json:"user_id" db:"user_id" gorm:"type:varchar(36);not null;index"`
	MovieID string    `json:"movie_id" db:"movie_id" gorm:"type:varchar(24);not null;index"` // MongoDB ObjectID as string
	AddedAt time.Time `json:"added_at" db:"added_at" gorm:"not null;default:now()"`
	Notes   string    `json:"notes" db:"notes" gorm:"type:text"`

	// GORM fields
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete support
}

func (WatchlistEntry) TableName() string {
	return "watchlist_entries"
}
