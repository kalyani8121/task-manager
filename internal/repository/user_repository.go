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
		INSERT INTO users (name, email, password)
		VALUES (:name, :email, :password)
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