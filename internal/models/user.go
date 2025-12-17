package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
// GORM model: automatically manages ID, created_at, updated_at, deleted_at
// Kept separate from database/HTTP representations for flexibility
type User struct {
	ID           string    `json:"id" gorm:"primaryKey;type:char(36)"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null;type:varchar(255)"`
	PasswordHash string    `json:"-" gorm:"not null;type:varchar(255)"`
	FirstName    string    `json:"first_name" gorm:"type:varchar(255)"`
	LastName     string    `json:"last_name" gorm:"type:varchar(255)"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name in database
// Prevents GORM from pluralizing (would use 'users' by default)
func (User) TableName() string {
	return "users"
}

// LoginRequest represents incoming login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents incoming registration request payload
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// AuthResponse represents successful authentication response
type AuthResponse struct {
	Token     string `json:"token"`
	User      *User  `json:"user"`
	ExpiresAt int64  `json:"expires_at"`
}
