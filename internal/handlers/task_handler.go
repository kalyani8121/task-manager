package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/kalyani8121/task-manager/internal/models"
	"github.com/kalyani8121/task-manager/internal/service"
)

type TaskHandler struct {
	service service.TaskService
	logger  *zap.Logger
}

func NewTaskHandler(svc service.TaskService, logger *zap.Logger) *TaskHandler {
	return &TaskHandler{service: svc, logger: logger}
}

// getUserID extracts the authenticated user's ID from the Gin context.
// This set by  AuthMiddleware.
func getUserID(c *gin.Context) string {
	return c.GetString("userID")
}

// CreateTask handles POST /api/v1/tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID := getUserID(c)
	var req models.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.CreateTask(userID, &req)
	if err != nil {
		h.logger.Error("Failed to create task", zap.String("userID", userID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Task created", zap.String("taskID", task.ID), zap.String("userID", userID))
	c.JSON(http.StatusCreated, task)
}

// GetTasks handles GET /api/v1/tasks
func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID := getUserID(c)

	tasks, err := h.service.GetTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// GetTask handles GET /api/v1/tasks/:id
func (h *TaskHandler) GetTask(c *gin.Context) {
	userID := getUserID(c)
	taskID := c.Param("id") // Reads :id from the URL

	task, err := h.service.GetTaskByID(taskID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask handles PUT /api/v1/tasks/:id
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID := getUserID(c)
	taskID := c.Param("id")

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.UpdateTask(taskID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Task updated", zap.String("taskID", taskID))
	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE /api/v1/tasks/:id
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID := getUserID(c)
	taskID := c.Param("id")

	if err := h.service.DeleteTask(taskID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Task deleted", zap.String("taskID", taskID))
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}