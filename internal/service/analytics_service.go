package service

import (
	"fmt"

	"github.com/kalyani8121/task-manager/internal/models"
	"github.com/kalyani8121/task-manager/internal/repository"
)

// AnalyticsService defines available analytics operations
type AnalyticsService interface {
	GetSummary(userID string) (*models.AnalyticsSummary, error)
	GetByStatus(userID string) ([]models.StatusCount, error)
	GetOverdue(userID string) (*models.OverdueResponse, error)
	GetPriorityQueue(userID string) (*models.PriorityQueueResponse, error)

}

type analyticsService struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsService(repo repository.AnalyticsRepository) AnalyticsService {
	return &analyticsService{repo: repo}
}

func (s *analyticsService) GetSummary(userID string) (*models.AnalyticsSummary, error) {
	summary, err := s.repo.GetSummary(userID)
	if err != nil {
		return nil, fmt.Errorf("get summary: %w", err)
	}
	return summary, nil
}

func (s *analyticsService) GetByStatus(userID string) ([]models.StatusCount, error) {
	counts, err := s.repo.GetByStatus(userID)
	if err != nil {
		return nil, fmt.Errorf("get by status: %w", err)
	}
	return counts, nil
}

func (s *analyticsService) GetOverdue(userID string) (*models.OverdueResponse, error) {
	tasks, err := s.repo.GetOverdue(userID)
	if err != nil {
		return nil, fmt.Errorf("get overdue: %w", err)
	}

	return &models.OverdueResponse{
		OverdueCount: len(tasks),
		Tasks:        tasks,
	}, nil
}

// GetPriorityQueue ranks tasks by urgency automatically
func (s *analyticsService) GetPriorityQueue(userID string) (*models.PriorityQueueResponse, error) {
	tasks, err := s.repo.GetPriorityQueue(userID)
	if err != nil {
		return nil, fmt.Errorf("get priority queue: %w", err)
	}

	// Calculate priority score for each task
	for i := range tasks {
		tasks[i].PriorityScore, tasks[i].Urgency = calculatePriority(tasks[i].DaysUntilDue)
		tasks[i].Rank = i+1
	}

	return &models.PriorityQueueResponse{
		TotalTasks: len(tasks),
		PriorityTasks: tasks,
	}, nil
}

func calculatePriority(daysUntilDue *int) (int, string) {
	if daysUntilDue == nil {
		return 0, "NO DEADLINE"
	}

	days := *daysUntilDue
	switch {
	case days < 0:
		return 100, "OVERDUE"
	case days == 0:
		return 95, "DUE TODAY"
	case days == 1:
		return 90, "DUE TOMORROW"
	case days <= 3:
		return 80, "URGENT"
	case days <= 7:
		return 60, "UPCOMING"
	case days <= 14:
		return 40, "NORMAL"
	default:
		return 20, "LOW PRIORITY"
	}
}