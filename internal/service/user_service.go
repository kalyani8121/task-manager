package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/kalyani8121/task-manager/internal/models"
	"github.com/kalyani8121/task-manager/internal/repository"
)

type UserService interface {
	Register(req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(req *models.LoginRequest) (*models.AuthResponse, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
	jwtExpiry int
}

func NewUserService(repo repository.UserRepository, jwtSecret string, jwtExpiry int) UserService {
	return &userService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *userService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if email already exists
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("checking existing user: %w", err)
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash the password — NEVER store plain text passwords.
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	// Create user in DB
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}
	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *userService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}

	// Use a generic error — don't tell attackers which field is wrong.
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

// generateToken creates a signed JWT token containing the user's ID.
func (s *userService) generateToken(userID string) (string, error) {
	// Claims are the "payload" of the JWT — what data it carries.
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(s.jwtExpiry) * time.Hour).Unix(),
		"iat":     time.Now().Unix(), // issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}