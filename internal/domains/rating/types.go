package rating

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MovieRating struct {
	ID        string         `json:"id" db:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID    string         `json:"user_id" db:"user_id" gorm:"type:varchar(36);not null;index"`
	MovieID   string         `json:"movie_id" db:"movie_id" gorm:"type:varchar(24);not null;index"` // MongoDB ObjectID as string
	Rating    float64        `json:"rating" db:"rating" gorm:"type:decimal(3,2);not null;check:rating >= 0 AND rating <= 5"`
	Review    string         `json:"review" db:"review" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete support
}

func (MovieRating) TableName() string {
	return "movie_ratings"
}

func (mr *MovieRating) BeforeCreate(tx *gorm.DB) error {
	// Validate rating is within range
	if mr.Rating < 0 || mr.Rating > 5 {
		return fmt.Errorf("rating must be between 0 and 5")
	}
	return nil
}
