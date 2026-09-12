package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/kalyani8121/task-manager/internal/handlers"
	"github.com/kalyani8121/task-manager/internal/middleware"
)

// SetupRoutes wires all URLs to their handlers.
func SetupRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	taskHandler *handlers.TaskHandler,
	analyticsHandler *handlers.AnalyticsHandler,
	jwtSecret string,
	logger *zap.Logger,
) {
	// Health check — used by Kubernetes to verify the pod is alive
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API versioning
	v1 := router.Group("/api/v1")
	{
		// Public routes — no token required
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// Protected routes — JWT middleware runs first
		tasks := v1.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware(jwtSecret))
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.GetTasks)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
		}

		analytics := v1.Group("/analytics")
		analytics.Use(middleware.AuthMiddleware(jwtSecret))
		{
			analytics.GET("/summary", analyticsHandler.GetSummary)
			analytics.GET("/by-status", analyticsHandler.GetByStatus)
			analytics.GET("/overdue", analyticsHandler.GetOverdue)
		}
	}
}
