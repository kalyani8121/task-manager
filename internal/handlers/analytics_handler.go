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

// GetSummary handles GET /api/v1/analytics/summary
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

// GetByStatus handles GET /api/v1/analytics/by-status
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

// GetOverdue handles GET /api/v1/analytics/overdue
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

// GetPriorityQueue handles GET /api/v1/tasks/priority-queue
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