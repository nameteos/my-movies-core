package user

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `json:"id" db:"id" gorm:"primaryKey;type:varchar(36)"`
	Username  string         `json:"username" db:"username" gorm:"type:varchar(100);not null;uniqueIndex"`
	Email     string         `json:"email" db:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	CreatedAt time.Time      `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete support
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Add any user validation logic here
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if u.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	return nil
}
