package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/dto"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/service"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/apperror"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/response"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	var query dto.DashboardStatsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	stats, err := h.dashboardService.GetStats(c.Request.Context(), userID, query.ProjectID)
	if err != nil {
		response.Error(c, err)
		return
	}

	tasksPerUser := make([]dto.TasksPerUserResponse, 0, len(stats.TasksPerUser))
	for _, item := range stats.TasksPerUser {
		tasksPerUser = append(tasksPerUser, dto.TasksPerUserResponse{
			UserID:   item.UserID,
			UserName: item.UserName,
			Count:    item.Count,
		})
	}

	response.Success(c, http.StatusOK, "dashboard stats fetched", dto.DashboardStatsResponse{
		ProjectID:     stats.ProjectID,
		TotalTasks:    stats.TotalTasks,
		OverdueTasks:  stats.OverdueTasks,
		TasksByStatus: stats.TasksByStatus,
		TasksPerUser:  tasksPerUser,
	})
}
