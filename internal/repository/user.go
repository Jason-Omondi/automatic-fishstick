package repository

import (
	"context"
	"errors"

	"github.com/Jason-Omondi/ecom/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserRepository handles all user-related database operations
// Repository pattern: abstracts data access logic with GORM
// GORM provides database-agnostic queries - switch MySQL↔PostgreSQL seamlessly
type UserRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewUserRepository(db *gorm.DB, log *zap.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

// CreateUser inserts a new user into the database
// Returns: newly created user with ID, or error if insert fails
// Why here: GORM handles database-specific SQL generation automatically
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	// GORM generates appropriate INSERT for MySQL or PostgreSQL
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.log.Error("Failed to create user", zap.String("email", user.Email), zap.Error(err))
		return err
	}

	r.log.Info("User created successfully", zap.String("email", user.Email), zap.String("id", user.ID))
	return nil
}

// GetUserByEmail retrieves a user from database by email
// Returns: user object if found, error if not found or query fails
// Why here: encapsulates query logic, GORM generates correct SQL for current DB
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	// GORM queries are database-agnostic
	// Same code works for MySQL, PostgreSQL, SQLite, etc.
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warn("User not found", zap.String("email", email))
			return nil, errors.New("user not found")
		}
		r.log.Error("Failed to fetch user", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user from database by ID
// Returns: user object if found, error if not found or query fails
// Why here: ID-based lookup common in auth flows after token validation
func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warn("User not found", zap.String("id", id))
			return nil, errors.New("user not found")
		}
		r.log.Error("Failed to fetch user", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return user, nil
}

// UpdateUser updates an existing user
// Returns: error if update fails
// Why here: provides abstraction for user updates across different databases
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		r.log.Error("Failed to update user", zap.String("id", user.ID), zap.Error(err))
		return err
	}

	r.log.Info("User updated successfully", zap.String("id", user.ID))
	return nil
}
