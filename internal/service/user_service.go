package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/kalyani8121/task-manager/internal/email"
	"github.com/kalyani8121/task-manager/internal/models"
	"github.com/kalyani8121/task-manager/internal/repository"
)

type UserService interface {
	Register(req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(req *models.LoginRequest) (*models.AuthResponse, error)
	VerifyEmail(token string) error
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
	jwtExpiry int
	mailer    *email.EmailSender
	appURL    string
}

func NewUserService(
	repo repository.UserRepository,
	jwtSecret string,
	jwtExpiry int,
	mailer *email.EmailSender,
	appURL string,
) UserService {
	return &userService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
		mailer:    mailer,
		appURL:    appURL,
	}
}

func isValidPassword(password string) error {
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' ||
			char == '#' || char == '$' ||
			char == '%' || char == '^' ||
			char == '&' || char == '*':
			hasSpecial = true
		}
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if !hasUpper {
		return errors.New("password must have at least 1 uppercase letter")
	}
	if !hasLower {
		return errors.New("password must have at least 1 lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must have at least 1 number")
	}
	if !hasSpecial {
		return errors.New("password must have at least 1 special character (!@#$%^&*)")
	}
	return nil
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *userService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Validate password strength
	if err := isValidPassword(req.Password); err != nil {
		return nil, err
	}

	// Check if email already exists
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("checking existing user: %w", err)
	}

	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	// Generate verification token
	verificationToken := generateToken()

	// Create user as NOT verified
	user := &models.User{
		Name:              req.Name,
		Email:             req.Email,
		Password:          string(hashed),
		IsVerified:        false,
		VerificationToken: verificationToken,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	// Send verification email
	verifyURL := fmt.Sprintf(
		"%s/api/v1/auth/verify?token=%s",
		s.appURL,
		verificationToken,
	)

	err = s.mailer.SendVerificationEmail(
		user.Email,
		user.Name,
		verifyURL,
	)

	if err != nil {
		return nil, fmt.Errorf("sending verification email: %w", err)
	}

	// User must verify email before logging in
	return nil, nil
}

func (s *userService) VerifyEmail(token string) error {
	user, err := s.repo.FindByVerificationToken(token)
	if err != nil {
		return fmt.Errorf("finding token: %w", err)
	}

	if user == nil {
		return errors.New("invalid or expired verification token")
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	return s.repo.MarkAsVerified(user.ID)
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

	// Check if email is verified
	if !user.IsVerified {
		return nil, errors.New("please verify your email before logging in")
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
