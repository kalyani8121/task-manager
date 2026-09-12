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