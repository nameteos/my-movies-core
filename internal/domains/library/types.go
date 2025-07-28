package library

import (
	"time"

	"gorm.io/gorm"
)

type WatchHistory struct {
	ID        string    `json:"id" db:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID    string    `json:"user_id" db:"user_id" gorm:"type:varchar(36);not null;index"`
	MovieID   string    `json:"movie_id" db:"movie_id" gorm:"type:varchar(24);not null;index"` // MongoDB ObjectID as string
	WatchedAt time.Time `json:"watched_at" db:"watched_at" gorm:"not null"`
	Duration  int       `json:"duration_watched" db:"duration_watched" gorm:"default:0"` // How long they watched in minutes

	// GORM fields
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete support
}

func (WatchHistory) TableName() string {
	return "watch_history"
}
