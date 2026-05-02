package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/apperror"
)

type successEnvelope struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type errorEnvelope struct {
	Error apperror.AppError `json:"error"`
}

func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, successEnvelope{
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, err error) {
	appErr := apperror.FromError(err)
	c.JSON(appErr.StatusCode, errorEnvelope{
		Error: *appErr,
	})
}

func NoContent(c *gin.Context, message string) {
	c.JSON(http.StatusOK, successEnvelope{
		Message: message,
	})
}
