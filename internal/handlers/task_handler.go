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

// CreateTask godoc
// @Summary      Create a task
// @Description  Create a new task for logged in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body models.CreateTaskRequest true "Create Task"
// @Success      201  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /tasks [post]
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

// GetTasks godoc
// @Summary      Get all tasks
// @Description  Get all tasks for logged in user
// @Tags         tasks
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /tasks [get]
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

// GetTask godoc
// @Summary      Get one task
// @Description  Get a specific task by ID
// @Tags         tasks
// @Produce      json
// @Param        id path string true "Task ID"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /tasks/{id} [get]
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

// UpdateTask godoc
// @Summary      Update a task
// @Description  Update task status or details
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Task ID"
// @Param        request body models.UpdateTaskRequest true "Update Task"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /tasks/{id} [put]
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

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Delete a task by ID
// @Tags         tasks
// @Produce      json
// @Param        id path string true "Task ID"
// @Success      200  {object}  map[string]string
// @Security     BearerAuth
// @Router       /tasks/{id} [delete]
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
