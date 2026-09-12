package models

import "time"

// TaskStatus is a custom string type — only these values are allowed.
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusCancelled  TaskStatus = "cancelled"
)

// Task maps to the `tasks` table.
type Task struct {
	ID          string     `db:"id"          json:"id"`
	UserID      string     `db:"user_id"     json:"user_id"`
	Title       string     `db:"title"       json:"title"`
	Description string     `db:"description" json:"description"`
	Status      TaskStatus `db:"status"      json:"status"`
	DueDate     *time.Time `db:"due_date"    json:"due_date"` // pointer = nullable
	CreatedAt   time.Time  `db:"created_at"  json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"  json:"updated_at"`
}

// CreateTaskRequest — what the client sends to create a task.
type CreateTaskRequest struct {
	Title       string     `json:"title"       binding:"required,min=1,max=255"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTaskRequest — what the client sends to update a task.
type UpdateTaskRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *TaskStatus `json:"status"`
	DueDate     *time.Time  `json:"due_date"`
}