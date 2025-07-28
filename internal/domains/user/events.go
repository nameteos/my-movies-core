package user

import (
	"event-driven-go/internal/shared"
)

const (
	UserRegisteredEventType = "user.user_registered"
	UserUpdatedEventType    = "user.user_updated"
	UserDeletedEventType    = "user.user_deleted"
)

type UserRegisteredEvent struct {
	shared.BaseEvent
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func init() {
	shared.RegisterEventType(UserRegisteredEventType, func() shared.Event {
		return &UserRegisteredEvent{}
	})
	shared.RegisterEventType(UserUpdatedEventType, func() shared.Event {
		return &UserUpdatedEvent{}
	})
	shared.RegisterEventType(UserDeletedEventType, func() shared.Event {
		return &UserDeletedEvent{}
	})
}

func (e UserRegisteredEvent) GetPayload() interface{} {
	return struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		UserID:   e.UserID,
		Username: e.Username,
		Email:    e.Email,
	}
}

func NewUserRegisteredEvent(userID, username, email string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		BaseEvent: shared.NewBaseEvent(UserRegisteredEventType),
		UserID:    userID,
		Username:  username,
		Email:     email,
	}
}

type UserUpdatedEvent struct {
	shared.BaseEvent
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (e UserUpdatedEvent) GetPayload() interface{} {
	return struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		UserID:   e.UserID,
		Username: e.Username,
		Email:    e.Email,
	}
}

func NewUserUpdatedEvent(userID, username, email string) *UserUpdatedEvent {
	return &UserUpdatedEvent{
		BaseEvent: shared.NewBaseEvent(UserUpdatedEventType),
		UserID:    userID,
		Username:  username,
		Email:     email,
	}
}

type UserDeletedEvent struct {
	shared.BaseEvent
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

func (e UserDeletedEvent) GetPayload() interface{} {
	return struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
	}{
		UserID:   e.UserID,
		Username: e.Username,
	}
}

func NewUserDeletedEvent(userID, username string) *UserDeletedEvent {
	return &UserDeletedEvent{
		BaseEvent: shared.NewBaseEvent(UserDeletedEventType),
		UserID:    userID,
		Username:  username,
	}
}
