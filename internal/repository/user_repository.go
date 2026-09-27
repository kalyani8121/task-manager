package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kalyani8121/task-manager/internal/models"
)

// UserRepository defines the contract (interface) for user DB operations.
type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	FindByVerificationToken(token string) (*models.User, error)
	MarkAsVerified(id string) error
}

// postgresUserRepository is the concrete implementation using PostgreSQL.
type postgresUserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

// Create inserts a new user into the database.
func (r *postgresUserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (name, email, password, is_verified, verification_token)
		VALUES (:name, :email, :password, :is_verified, :verification_token)
		RETURNING id, created_at, updated_at
	`
	// NamedQuery uses struct field names (via db tags) as SQL parameters.
	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	}
	return nil
}

// FindByEmail retrieves a user by their email address.
func (r *postgresUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, password, created_at, updated_at
              FROM users WHERE email = $1`

	err := r.db.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found — return nil, nil (not an error)
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &user, nil
}

// FindByID retrieves a user by their UUID.
func (r *postgresUserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, created_at, updated_at
              FROM users WHERE id = $1`

	err := r.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

func (r *postgresUserRepository) FindByVerificationToken(token string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, password, 
			  is_verified, verification_token,
			  created_at, updated_at
              FROM users WHERE verification_token = $1`

	err := r.db.Get(&user, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find by token: %w", err)
	}
	return &user, nil
}

func (r *postgresUserRepository) MarkAsVerified(id string) error {
	query := `UPDATE users 
			  SET is_verified = TRUE, 
			      verification_token = NULL,
			      updated_at = NOW()
			  WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("mark verified: %w", err)
	}
	return nil
}
