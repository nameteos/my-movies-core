package user

import (
	"context"
	"fmt"

	"event-driven-go/internal/shared"
)

type Service struct {
	repository RepositoryInterface
	eventBus   *shared.EventBus
}

func NewService(repository RepositoryInterface, eventBus *shared.EventBus) *Service {
	return &Service{
		repository: repository,
		eventBus:   eventBus,
	}
}

func (s *Service) RegisterUser(ctx context.Context, username, email string) (*User, error) {
	// Validate input
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	// Check if user already exists
	exists, err := s.repository.UserExists(ctx, username, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user with username '%s' or email '%s' already exists", username, email)
	}

	// Create user in repository
	user, err := s.repository.CreateUser(ctx, username, email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Publish user registered event
	event := NewUserRegisteredEvent(user.ID, user.Username, user.Email)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation since user was created
		fmt.Printf("Warning: failed to publish user registered event: %v\n", err)
	}

	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	return s.repository.GetUserByID(ctx, id)
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	return s.repository.GetUserByUsername(ctx, username)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	return s.repository.GetUserByEmail(ctx, email)
}

func (s *Service) UpdateUser(ctx context.Context, user *User) (*User, error) {
	// Validate input
	if user == nil {
		return nil, fmt.Errorf("user cannot be nil")
	}
	if user.ID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	// Update user in repository
	updatedUser, err := s.repository.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Publish user updated event
	event := NewUserUpdatedEvent(updatedUser.ID, updatedUser.Username, updatedUser.Email)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to publish user updated event: %v\n", err)
	}

	return updatedUser, nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	// Get user first to get username for event
	user, err := s.repository.GetUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user before deletion: %w", err)
	}

	// Delete user from repository
	if err := s.repository.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Publish user deleted event
	event := NewUserDeletedEvent(id, user.Username)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation since user was deleted
		fmt.Printf("Warning: failed to publish user deleted event: %v\n", err)
	}

	return nil
}

func (s *Service) ListUsers(ctx context.Context, limit, offset int) ([]*User, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.repository.ListUsers(ctx, limit, offset)
}

func (s *Service) GetActiveUserCount(ctx context.Context) (int64, error) {
	return s.repository.GetActiveUserCount(ctx)
}

func (s *Service) GetRecentUsers(ctx context.Context, limit int) ([]*User, error) {
	if limit <= 0 {
		limit = 10
	}

	return s.repository.GetRecentUsers(ctx, limit)
}
