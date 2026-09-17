package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kalyani8121/task-manager/internal/models"
)

// AnalyticsRepository defines what analytics
type AnalyticsRepository interface {
	GetSummary(userID string) (*models.AnalyticsSummary, error)
	GetByStatus(userID string) ([]models.StatusCount, error)
	GetOverdue(userID string) ([]models.OverdueTask, error)
	GetPriorityQueue(userID string) ([]models.PriorityTask,error)
}

type postgresAnalyticsRepository struct {
	db *sqlx.DB
}

func NewAnalyticsRepository(db *sqlx.DB) AnalyticsRepository {
	return &postgresAnalyticsRepository{db: db}
}

// GetSummary gets overall task statistics
func (r *postgresAnalyticsRepository) GetSummary(userID string) (*models.AnalyticsSummary, error) {

	// Count total tasks
	var total int
	err := r.db.Get(&total,
		`SELECT COUNT(*) FROM tasks WHERE user_id = $1`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("count total: %w", err)
	}

	// Count completed tasks
	var completed int
	err = r.db.Get(&completed,
		`SELECT COUNT(*) FROM tasks
		 WHERE user_id = $1 AND status = 'completed'`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("count completed: %w", err)
	}

	// Count pending tasks
	var pending int
	err = r.db.Get(&pending,
		`SELECT COUNT(*) FROM tasks
		 WHERE user_id = $1 AND status = 'pending'`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("count pending: %w", err)
	}

	// Count in_progress tasks
	var inProgress int
	err = r.db.Get(&inProgress,
		`SELECT COUNT(*) FROM tasks
		 WHERE user_id = $1 AND status = 'in_progress'`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("count in_progress: %w", err)
	}

	// Count cancelled tasks
	var cancelled int
	err = r.db.Get(&cancelled,
		`SELECT COUNT(*) FROM tasks
		 WHERE user_id = $1 AND status = 'cancelled'`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("count cancelled: %w", err)
	}

	// Count overdue tasks
	var overdue int
	err = r.db.Get(&overdue,
		`SELECT COUNT(*) FROM tasks
		 WHERE user_id = $1
		 AND due_date < NOW()
		 AND status NOT IN ('completed', 'cancelled')`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("count overdue: %w", err)
	}

	// Calculate completion rate
	var completionRate float64
	if total > 0 {
		completionRate = float64(completed) / float64(total) * 100
	}

	return &models.AnalyticsSummary{
		TotalTasks:     total,
		Completed:      completed,
		Pending:        pending,
		InProgress:     inProgress,
		Cancelled:      cancelled,
		Overdue:        overdue,
		CompletedPercent: completionRate,
	}, nil
}

// GetByStatus counts tasks grouped by status
func (r *postgresAnalyticsRepository) GetByStatus(userID string) ([]models.StatusCount, error) {
	var counts []models.StatusCount

	query := `
		SELECT status, COUNT(*) as count
		FROM tasks
		WHERE user_id = $1
		GROUP BY status
		ORDER BY count DESC
	`

	err := r.db.Select(&counts, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get by status: %w", err)
	}

	return counts, nil
}

// GetOverdue gets all overdue tasks
func (r *postgresAnalyticsRepository) GetOverdue(userID string) ([]models.OverdueTask, error) {
	var tasks []models.OverdueTask

	query := `
		SELECT
			id,
			title,
			TO_CHAR(due_date, 'YYYY-MM-DD') as due_date,
			EXTRACT(DAY FROM NOW() - due_date)::int as days_overdue,
			status::text as status
		FROM tasks
		WHERE user_id = $1
		AND due_date < NOW()
		AND status NOT IN ('completed', 'cancelled')
		ORDER BY due_date ASC
	`

	err := r.db.Select(&tasks, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get overdue: %w", err)
	}

	return tasks, nil
}

// GetPriorityQueue gets tasks ordered by priority and due date
func (r *postgresAnalyticsRepository) GetPriorityQueue(userID string) ([]models.PriorityTask, error) {
	var tasks []models.PriorityTask

	query := `
		SELECT
			id,
			title,
			status::text as status,
			CASE
				WHEN due_date IS NULL THEN NULL
				ELSE TO_CHAR(due_date, 'YYYY-MM-DD')
			END as due_date,
			CASE
				WHEN due_date IS NULL THEN NULL
				ELSE EXTRACT(DAY FROM due_date - NOW())::int
			END as days_until_due
		FROM tasks
		WHERE user_id = $1
		AND status NOT IN ('completed', 'cancelled')
		ORDER BY
			CASE
				WHEN due_date IS NULL THEN 2
				ELSE 1
			END,
			due_date ASC
	`

	err := r.db.Select(&tasks, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get priority queue: %w", err)
	}

	return tasks, nil
}