package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/dto"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/service"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/apperror"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request dto.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	user, token, err := h.authService.Register(c.Request.Context(), service.RegisterInput{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "registration successful", dto.NewAuthResponse(user, token))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	user, token, err := h.authService.Login(c.Request.Context(), service.LoginInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "login successful", dto.NewAuthResponse(user, token))
}
