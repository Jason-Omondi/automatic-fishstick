package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Jason-Omondi/ecomgo/internal/models"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	// Service layer handles business logic
	// Handler only coordinates HTTP request/response and delegates to service
	service *UserService
	log     *zap.Logger
}

func NewHandler(service *UserService, log *zap.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// RegisterRoutes registers user-related routes to the given router
// Routes define HTTP endpoints and map them to handler methods
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/users/{id}", h.handleGetUser).Methods("GET")
}

// handleRegister handles POST /api/v1/register
// @Summary Register new user
// @Description Creates a new user account and returns auth token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {string} string "Invalid request or user already exists"
// @Failure 500 {string} string "Internal server error"
// @Router /register [post]
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Register endpoint called")

	var req models.RegisterRequest
	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("Invalid register request", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Call service to handle registration logic
	authResp, err := h.service.Register(context.Background(), &req)
	if err != nil {
		h.log.Error("Registration failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authResp)
}

// handleLogin handles POST /api/v1/login
// @Summary Login user
// @Description Authenticates user and returns auth token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Login endpoint called")

	var req models.LoginRequest
	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("Invalid login request", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Call service to handle login logic
	authResp, err := h.service.Login(context.Background(), &req)
	if err != nil {
		h.log.Warn("Login failed", zap.Error(err))
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Return successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResp)
}

// handleGetUser handles GET /api/v1/users/{id}
// @Summary Get user by ID
// @Description Retrieves user data by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{id} [get]
func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	vars := mux.Vars(r)
	userID := vars["id"]

	h.log.Info("Get user endpoint called", zap.String("id", userID))

	// Call service to fetch user
	user, err := h.service.GetUserByID(context.Background(), userID)
	if err != nil {
		h.log.Warn("User not found", zap.String("id", userID), zap.Error(err))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Return successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
