package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/kalyani8121/task-manager/internal/service"
)

type AnalyticsHandler struct {
	service service.AnalyticsService
	logger  *zap.Logger
}

func NewAnalyticsHandler(svc service.AnalyticsService, logger *zap.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{service: svc, logger: logger}
}

// GetSummary godoc
// @Summary      Get analytics summary
// @Description  Get overall task statistics
// @Tags         analytics
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /analytics/summary [get]
func (h *AnalyticsHandler) GetSummary(c *gin.Context) {
	userID := getUserID(c)

	summary, err := h.service.GetSummary(userID)
	if err != nil {
		h.logger.Error("Failed to get summary",
			zap.String("userID", userID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Analytics summary fetched",
		zap.String("userID", userID))
	c.JSON(http.StatusOK, summary)
}

// GetByStatus godoc
// @Summary      Get tasks by status
// @Description  Get count of tasks grouped by status
// @Tags         analytics
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /analytics/by-status [get]
func (h *AnalyticsHandler) GetByStatus(c *gin.Context) {
	userID := getUserID(c)

	counts, err := h.service.GetByStatus(userID)
	if err != nil {
		h.logger.Error("Failed to get by status",
			zap.String("userID", userID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_breakdown": counts,
	})
}

// GetOverdue godoc
// @Summary      Get overdue tasks
// @Description  Get all tasks past their due date
// @Tags         analytics
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /analytics/overdue [get]
func (h *AnalyticsHandler) GetOverdue(c *gin.Context) {
	userID := getUserID(c)

	overdue, err := h.service.GetOverdue(userID)
	if err != nil {
		h.logger.Error("Failed to get overdue",
			zap.String("userID", userID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, overdue)
}

// GetPriorityQueue godoc
// @Summary      Get priority queue
// @Description  Get tasks ranked by urgency automatically
// @Tags         analytics
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /tasks/priority-queue [get]
func (h *AnalyticsHandler) GetPriorityQueue(c *gin.Context) {
	userID := getUserID(c)

	result, err := h.service.GetPriorityQueue(userID)
	if err != nil {
		h.logger.Error("Failed to get priority queue",
			zap.String("userID", userID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Priority queue fetched",
		zap.String("userID", userID),
		zap.Int("total",result.TotalTasks))
	c.JSON(http.StatusOK, result)
}