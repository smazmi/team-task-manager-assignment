package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/middleware"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/apperror"
)

func currentUserID(c *gin.Context) (uint, error) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		return 0, apperror.Unauthorized("authenticated user not found in request context")
	}

	return userID, nil
}

func parseUintParam(c *gin.Context, key string) (uint, error) {
	rawValue := c.Param(key)
	value, err := strconv.ParseUint(rawValue, 10, 64)
	if err != nil {
		return 0, apperror.BadRequest("invalid path parameter: " + key)
	}

	return uint(value), nil
}
