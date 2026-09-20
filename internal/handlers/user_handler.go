package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/kalyani8121/task-manager/internal/models"
	"github.com/kalyani8121/task-manager/internal/service"
)

type UserHandler struct {
	service service.UserService
	logger  *zap.Logger
}

func NewUserHandler(svc service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: svc, logger: logger}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.RegisterRequest true "Register Request"
// @Success      201  {object}  models.AuthResponse
// @Failure      400  {object}  map[string]string
// @Router       /auth/register [post]

// Register handles POST /api/v1/auth/register
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	// ShouldBindJSON parses the JSON body AND runs validation (binding tags).
	// If name is missing, or email is invalid, this returns an error automatically.
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Register(&req)
	if err != nil {
		h.logger.Warn("Registration failed", zap.String("email", req.Email), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("User registered", zap.String("email", req.Email))
	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary      Login user
// @Description  Login with email and password to get JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Login Request"
// @Success      200  {object}  models.AuthResponse
// @Failure      401  {object}  map[string]string
// @Router       /auth/login [post]

// Login handles POST /api/v1/auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Login(&req)
	if err != nil {
		// Use 401 for auth failures, not 400
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("User logged in", zap.String("email", req.Email))
	c.JSON(http.StatusOK, resp)
}