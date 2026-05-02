package service

import (
	"context"
	"errors"
	"time"

	"github.com/smazmi/team-task-manager-assignment/backend/internal/models"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/repository"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/apperror"
	"gorm.io/gorm"
)

type DashboardTasksPerUser struct {
	UserID   uint
	UserName string
	Count    int64
}

type DashboardStats struct {
	ProjectID     uint
	TotalTasks    int64
	OverdueTasks  int64
	TasksByStatus map[string]int64
	TasksPerUser  []DashboardTasksPerUser
}

type DashboardService struct {
	projects repository.ProjectRepository
	tasks    repository.TaskRepository
}

func NewDashboardService(projects repository.ProjectRepository, tasks repository.TaskRepository) *DashboardService {
	return &DashboardService{
		projects: projects,
		tasks:    tasks,
	}
}

func (s *DashboardService) GetStats(ctx context.Context, userID, projectID uint) (*DashboardStats, error) {
	member, err := s.projects.GetMembership(ctx, projectID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Forbidden("you do not belong to this project")
		}
		return nil, apperror.Internal("failed to verify project membership")
	}

	var assigneeFilter *uint
	if member.Role != models.ProjectRoleAdmin {
		assigneeFilter = &userID
	}

	totalTasks, err := s.tasks.CountByProjectAndAssignee(ctx, projectID, assigneeFilter)
	if err != nil {
		return nil, apperror.Internal("failed to count tasks")
	}

	statusCounts, err := s.tasks.GetStatusCounts(ctx, projectID, assigneeFilter)
	if err != nil {
		return nil, apperror.Internal("failed to compute task status breakdown")
	}

	tasksPerUserRows, err := s.tasks.GetTasksPerUser(ctx, projectID, assigneeFilter)
	if err != nil {
		return nil, apperror.Internal("failed to compute assignee breakdown")
	}

	overdueTasks, err := s.tasks.CountOverdue(ctx, projectID, assigneeFilter, time.Now())
	if err != nil {
		return nil, apperror.Internal("failed to count overdue tasks")
	}

	tasksByStatus := map[string]int64{
		string(models.TaskStatusTodo):       0,
		string(models.TaskStatusInProgress): 0,
		string(models.TaskStatusDone):       0,
	}
	for status, count := range statusCounts {
		tasksByStatus[string(status)] = count
	}

	tasksPerUser := make([]DashboardTasksPerUser, 0, len(tasksPerUserRows))
	for _, row := range tasksPerUserRows {
		tasksPerUser = append(tasksPerUser, DashboardTasksPerUser{
			UserID:   row.UserID,
			UserName: row.UserName,
			Count:    row.Count,
		})
	}

	return &DashboardStats{
		ProjectID:     projectID,
		TotalTasks:    totalTasks,
		OverdueTasks:  overdueTasks,
		TasksByStatus: tasksByStatus,
		TasksPerUser:  tasksPerUser,
	}, nil
}
