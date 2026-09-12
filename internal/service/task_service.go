package service

import (
	"errors"
	"fmt"

	"github.com/kalyani8121/task-manager/internal/models"
	"github.com/kalyani8121/task-manager/internal/repository"
)

type TaskService interface {
	CreateTask(userID string, req *models.CreateTaskRequest) (*models.Task, error)
	GetTasks(userID string) ([]models.Task, error)
	GetTaskByID(id, userID string) (*models.Task, error)
	UpdateTask(id, userID string, req *models.UpdateTaskRequest) (*models.Task, error)
	DeleteTask(id, userID string) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) CreateTask(userID string, req *models.CreateTaskRequest) (*models.Task, error) {
	// Default status to "pending" if not provided
	status := req.Status
	if status == "" {
		status = models.StatusPending
	}

	// Validate the status value
	if !isValidStatus(status) {
		return nil, errors.New("invalid task status")
	}

	task := &models.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		DueDate:     req.DueDate,
	}

	if err := s.repo.Create(task); err != nil {
		return nil, fmt.Errorf("creating task: %w", err)
	}
	return task, nil
}

func (s *taskService) GetTasks(userID string) ([]models.Task, error) {
	tasks, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("fetching tasks: %w", err)
	}
	return tasks, nil
}

func (s *taskService) GetTaskByID(id, userID string) (*models.Task, error) {
	task, err := s.repo.FindByID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching task: %w", err)
	}
	if task == nil {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (s *taskService) UpdateTask(id, userID string, req *models.UpdateTaskRequest) (*models.Task, error) {
	// First fetch the existing task
	task, err := s.repo.FindByID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching task: %w", err)
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	// Partial update: only overwrite fields the client actually sent.
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		if !isValidStatus(*req.Status) {
			return nil, errors.New("invalid task status")
		}
		task.Status = *req.Status
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := s.repo.Update(task); err != nil {
		return nil, fmt.Errorf("updating task: %w", err)
	}
	return task, nil
}

func (s *taskService) DeleteTask(id, userID string) error {
	if err := s.repo.Delete(id, userID); err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}
	return nil
}

func isValidStatus(s models.TaskStatus) bool {
	switch s {
	case models.StatusPending, models.StatusInProgress, models.StatusCompleted, models.StatusCancelled:
		return true
	}
	return false
}