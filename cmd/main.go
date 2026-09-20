// @title           Task Management System API
// @version         1.0
// @description     Production-style Task Management REST API with JWT authentication, analytics and smart priority queue.

// @contact.name   Kalyani
// @contact.email  kalyanikuntumalla3@gmail.com

// @license.name  MIT

// @host      task-manager-production-4677.up.railway.app
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/kalyani8121/task-manager/internal/config"
	"github.com/kalyani8121/task-manager/internal/handlers"
	"github.com/kalyani8121/task-manager/internal/repository"
	"github.com/kalyani8121/task-manager/internal/routes"
	"github.com/kalyani8121/task-manager/internal/service"
	"github.com/kalyani8121/task-manager/internal/email"
	"github.com/kalyani8121/task-manager/internal/scheduler"

	_ "github.com/lib/pq"
)

func main() {
	// 1. Load configuration 
	cfg := config.LoadConfig()

	// 2. Initialize structured logger
	// In production, use zap.NewProduction() for JSON logs (Kibana-friendly)
	// In development, use zap.NewDevelopment() for human-readable logs
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync() // Flushes buffered log entries

	// 3. Connect to database 
	var db *sqlx.DB
	var connectErr error
	for i := 0; i < 10; i++ {
		db, connectErr = sqlx.Open("postgres", cfg.DBURL)
		if connectErr == nil {
			if connectErr = db.Ping(); connectErr == nil {
				break
			}
		}
		logger.Warn("Database not ready yet, retrying", zap.Int("attempt", i+1), zap.Error(connectErr))
		time.Sleep(2 * time.Second)
	}
	if connectErr != nil {
		logger.Fatal("Failed to connect to database after retries", zap.Error(connectErr))
	}
	defer db.Close()

	// 4. Wire up the layers (Dependency Injection)
	// Repository → Service → Handler
	// Each layer only knows about the layer below it (via interfaces).
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)


	userService := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	taskService := service.NewTaskService(taskRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)

	userHandler := handlers.NewUserHandler(userService, logger)
	taskHandler := handlers.NewTaskHandler(taskService, logger)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, logger)

	//Setup email sender
	mailer := email.NewEmailSender(
		cfg.ResendAPIKey,
        cfg.SMTPEmail,
	)
		// Start email scheduler (runs every day 9AM)
		emailScheduler := scheduler.NewScheduler(db, mailer)
	    emailScheduler.Start()
	

	// Create Gin router
	router := gin.New()
	router.RedirectTrailingSlash = false

	// Gin middleware: logging and crash recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery()) // Catches panics and returns 500 instead of crashing

	// 6. Register routes
	routes.SetupRoutes(router, userHandler, taskHandler, analyticsHandler, cfg.JWTSecret, logger)

	// 7. Start server 
	logger.Info("Server starting", zap.String("port", cfg.AppPort))
	if err := router.Run(":" + cfg.AppPort); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
