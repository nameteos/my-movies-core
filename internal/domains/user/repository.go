package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	CreateUser(ctx context.Context, username, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
	UserExists(ctx context.Context, username, email string) (bool, error)
	GetActiveUserCount(ctx context.Context) (int64, error)
	GetRecentUsers(ctx context.Context, limit int) ([]*User, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, username, email string) (*User, error) {
	user := &User{
		ID:       uuid.New().String(),
		Username: username,
		Email:    email,
	}

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create user: %w", result.Error)
	}

	return user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User

	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User

	result := r.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *User) (*User, error) {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update user: %w", result.Error)
	}

	return user, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *Repository) ListUsers(ctx context.Context, limit, offset int) ([]*User, error) {
	var users []*User

	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list users: %w", result.Error)
	}

	return users, nil
}

func (r *Repository) UserExists(ctx context.Context, username, email string) (bool, error) {
	var count int64

	result := r.db.WithContext(ctx).
		Model(&User{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check user existence: %w", result.Error)
	}

	return count > 0, nil
}

func (r *Repository) GetActiveUserCount(ctx context.Context) (int64, error) {
	var count int64

	result := r.db.WithContext(ctx).Model(&User{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to get active user count: %w", result.Error)
	}

	return count, nil
}

func (r *Repository) GetRecentUsers(ctx context.Context, limit int) ([]*User, error) {
	var users []*User

	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get recent users: %w", result.Error)
	}

	return users, nil
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&User{})
}
