package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/Jason-Omondi/ecomgo/internal/config"
	"github.com/Jason-Omondi/ecomgo/internal/models"
	"github.com/Jason-Omondi/ecomgo/internal/repository"
	"go.uber.org/zap"
)

// UserService implements business logic for user operations
// Service layer: coordinates between HTTP handlers and data repositories
// Config is injected once and reused for all operations
type UserService struct {
	userRepo *repository.UserRepository
	log      *zap.Logger
	config   *config.Config // Store config for Keycloak, external services, etc.
}

func NewUserService(userRepo *repository.UserRepository,
	log *zap.Logger, cfg *config.Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		log:      log,
		config:   cfg,
	}
}

// Register creates a new user account
// Uses config for validation rules, token expiry settings, etc.
func (s *UserService) Register(ctx context.Context,
	req *models.RegisterRequest) (*models.AuthResponse, error) {
	s.log.Info("Registering new user",
		zap.String("email", req.Email))

	// Check if user already exists
	existingUser, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		s.log.Warn("Registration failed: user already exists",
			zap.String("email", req.Email))
		return nil, errors.New("user already exists")
	}

	// Hash password - use bcrypt in production for security
	// SHA256 used here for demo; replace with golang.org/x/crypto/bcrypt for production
	hashedPassword := s.hashPassword(req.Password)

	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	// Create user in database
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		s.log.Error("Failed to create user",
			zap.String("email", req.Email), zap.Error(err))
		return nil, err
	}

	// Use config.Keycloak settings when generating real OAuth tokens
	token := s.generateToken(user.ID)

	s.log.Info("User registered successfully", zap.String("email", user.Email))

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// Login authenticates user
// Uses config for token expiry, Keycloak realm, etc.
func (s *UserService) Login(ctx context.Context,
	req *models.LoginRequest) (*models.AuthResponse, error) {
	s.log.Info("User login attempt", zap.String("email", req.Email))

	// Fetch user by email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.log.Warn("Login failed: user not found",
			zap.String("email", req.Email))
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if !s.verifyPassword(req.Password, user.PasswordHash) {
		s.log.Warn("Login failed: invalid password",
			zap.String("email", req.Email))
		return nil, errors.New("invalid credentials")
	}

	// Access config.Keycloak.URL, config.Keycloak.Realm etc. here
	token := s.generateToken(user.ID)

	s.log.Info("User logged in successfully",
		zap.String("email", user.Email))

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetUserByID retrieves user data by ID
// Returns: user object if found, error if not found
// Why here: delegates to repository after validating context
func (s *UserService) GetUserByID(ctx context.Context,
	id string) (*models.User, error) {
	s.log.Info("Fetching user data", zap.String("id", id))

	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		s.log.Warn("Failed to fetch user", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return user, nil
}

// hashPassword hashes password using SHA256
// For production: use golang.org/x/crypto/bcrypt instead
// Returns: hex-encoded hash string
func (s *UserService) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// verifyPassword compares plain password with stored hash
// Returns: true if password matches hash, false otherwise
func (s *UserService) verifyPassword(password, hash string) bool {
	return s.hashPassword(password) == hash
}

// generateToken creates a simple JWT-like token (implement real JWT with Keycloak in production)
// For production: integrate with Keycloak OAuth2 token endpoint
// Returns: token string
func (s *UserService) generateToken(userID string) string {
	// TODO: Use s.config.Keycloak to call real OAuth2 token endpoint
	// Example: s.config.Keycloak.URL, s.config.Keycloak.ClientID
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." + userID
}
