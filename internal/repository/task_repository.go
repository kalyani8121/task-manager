package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kalyani8121/task-manager/internal/models"
)

type TaskRepository interface {
	Create(task *models.Task) error
	FindByUserID(userID string) ([]models.Task, error)
	FindByID(id, userID string) (*models.Task, error)
	Update(task *models.Task) error
	Delete(id, userID string) error
}

type postgresTaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) TaskRepository {
	return &postgresTaskRepository{db: db}
}

func (r *postgresTaskRepository) Create(task *models.Task) error {
	query := `
		INSERT INTO tasks (user_id, title, description, status, due_date)
		VALUES (:user_id, :title, :description, :status, :due_date)
		RETURNING id, created_at, updated_at
	`
	rows, err := r.db.NamedQuery(query, task)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	}
	return nil
}

// FindByUserID gets all tasks for a specific user.
func (r *postgresTaskRepository) FindByUserID(userID string) ([]models.Task, error) {
	var tasks []models.Task
	query := `
		SELECT id, user_id, title, description, status, due_date, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	err := r.db.Select(&tasks, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find tasks by user: %w", err)
	}
	return tasks, nil
}

// FindByID gets one task — but ONLY if it belongs to the requesting user.
func (r *postgresTaskRepository) FindByID(id, userID string) (*models.Task, error) {
	var task models.Task
	query := `
		SELECT id, user_id, title, description, status, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2
	`
	err := r.db.Get(&task, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find task by id: %w", err)
	}
	return &task, nil
}

func (r *postgresTaskRepository) Update(task *models.Task) error {
	query := `
		UPDATE tasks
		SET title = :title,
		    description = :description,
		    status = :status,
		    due_date = :due_date,
		    updated_at = NOW()
		WHERE id = :id AND user_id = :user_id
	`
	result, err := r.db.NamedExec(query, task)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("task not found or not owned by user")
	}
	return nil
}

func (r *postgresTaskRepository) Delete(id, userID string) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("task not found or not owned by user")
	}
	return nil
}