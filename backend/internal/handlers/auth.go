package handlers

import (
	"biblia-am-pm/internal/middleware"
	"biblia-am-pm/internal/repository"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userRepo: repository.NewUserRepository(),
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Check if user already exists
	existingUser, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to check user")
		return
	}

	if existingUser != nil {
		middleware.RespondError(w, http.StatusConflict, "User already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create user
	user, err := h.userRepo.CreateUser(req.Email, string(hashedPassword))
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	middleware.RespondJSON(w, http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Get user by email
	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	if user == nil {
		middleware.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		middleware.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Remove password from response
	user.Password = ""

	middleware.RespondJSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

func generateToken(userID int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key"
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

