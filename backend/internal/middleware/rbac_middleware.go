package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/models"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/repository"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/apperror"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/response"
	"gorm.io/gorm"
)

type RBACMiddleware struct {
	projects repository.ProjectRepository
	tasks    repository.TaskRepository
}

func NewRBACMiddleware(projects repository.ProjectRepository, tasks repository.TaskRepository) *RBACMiddleware {
	return &RBACMiddleware{
		projects: projects,
		tasks:    tasks,
	}
}

func (m *RBACMiddleware) RequireProjectMember() gin.HandlerFunc {
	return m.RequireProjectRoles(models.ProjectRoleAdmin, models.ProjectRoleMember)
}

func (m *RBACMiddleware) RequireProjectAdmin() gin.HandlerFunc {
	return m.RequireProjectRoles(models.ProjectRoleAdmin)
}

func (m *RBACMiddleware) RequireProjectRoles(roles ...models.ProjectRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := GetCurrentUserID(c)
		if !ok {
			response.Error(c, apperror.Unauthorized("authenticated user not found in request context"))
			c.Abort()
			return
		}

		projectID, err := resolveProjectID(c)
		if err != nil {
			response.Error(c, apperror.BadRequest(err.Error()))
			c.Abort()
			return
		}

		allowed, err := m.projects.UserHasRole(c.Request.Context(), projectID, userID, roles...)
		if err != nil {
			response.Error(c, apperror.Internal("failed to verify project access"))
			c.Abort()
			return
		}

		if !allowed {
			response.Error(c, apperror.Forbidden("you do not have permission to access this project"))
			c.Abort()
			return
		}

		c.Set("project_id", projectID)
		c.Next()
	}
}

func (m *RBACMiddleware) RequireTaskAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		task, userID, ok := m.resolveTaskAndUser(c)
		if !ok {
			return
		}

		allowed, err := m.projects.UserHasRole(c.Request.Context(), task.ProjectID, userID, models.ProjectRoleAdmin)
		if err != nil {
			response.Error(c, apperror.Internal("failed to verify task admin access"))
			c.Abort()
			return
		}

		if !allowed {
			response.Error(c, apperror.Forbidden("admin access is required for this task"))
			c.Abort()
			return
		}

		c.Set("project_id", task.ProjectID)
		c.Set("task_id", task.ID)
		c.Next()
	}
}

func (m *RBACMiddleware) RequireTaskActor() gin.HandlerFunc {
	return func(c *gin.Context) {
		task, userID, ok := m.resolveTaskAndUser(c)
		if !ok {
			return
		}

		isAdmin, err := m.projects.UserHasRole(c.Request.Context(), task.ProjectID, userID, models.ProjectRoleAdmin)
		if err != nil {
			response.Error(c, apperror.Internal("failed to verify task permissions"))
			c.Abort()
			return
		}

		if isAdmin || (task.AssigneeID != nil && *task.AssigneeID == userID) {
			c.Set("project_id", task.ProjectID)
			c.Set("task_id", task.ID)
			c.Next()
			return
		}

		response.Error(c, apperror.Forbidden("you can only access tasks assigned to you"))
		c.Abort()
	}
}

func (m *RBACMiddleware) resolveTaskAndUser(c *gin.Context) (*models.Task, uint, bool) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("authenticated user not found in request context"))
		c.Abort()
		return nil, 0, false
	}

	taskIDValue := c.Param("taskId")
	taskID, err := strconv.ParseUint(taskIDValue, 10, 64)
	if err != nil {
		response.Error(c, apperror.BadRequest("invalid task id"))
		c.Abort()
		return nil, 0, false
	}

	task, err := m.tasks.GetByID(c.Request.Context(), uint(taskID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, apperror.NotFound("task not found"))
		} else {
			response.Error(c, apperror.Internal("failed to load task"))
		}
		c.Abort()
		return nil, 0, false
	}

	return task, userID, true
}

func resolveProjectID(c *gin.Context) (uint, error) {
	for _, key := range []string{"projectId", "project_id"} {
		if value := c.Param(key); value != "" {
			return parseUint(value)
		}
	}

	if value := c.Query("project_id"); value != "" {
		return parseUint(value)
	}

	if c.Request.Body == nil {
		return 0, fmt.Errorf("project_id is required")
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read request body")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if len(bytes.TrimSpace(bodyBytes)) == 0 {
		return 0, fmt.Errorf("project_id is required")
	}

	var payload map[string]any
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return 0, fmt.Errorf("project_id is required")
	}

	value, ok := payload["project_id"]
	if !ok {
		return 0, fmt.Errorf("project_id is required")
	}

	projectID, ok := toUint(value)
	if !ok {
		return 0, fmt.Errorf("project_id must be a positive integer")
	}

	return projectID, nil
}

func parseUint(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("project_id must be a positive integer")
	}

	return uint(parsed), nil
}

func toUint(value any) (uint, bool) {
	switch typed := value.(type) {
	case float64:
		if typed < 1 {
			return 0, false
		}
		return uint(typed), true
	case string:
		parsed, err := strconv.ParseUint(typed, 10, 64)
		if err != nil {
			return 0, false
		}
		return uint(parsed), true
	default:
		return 0, false
	}
}
