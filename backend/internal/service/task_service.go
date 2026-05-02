package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/smazmi/team-task-manager-assignment/backend/internal/models"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/repository"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/apperror"
	"gorm.io/gorm"
)

type CreateTaskInput struct {
	ProjectID   uint
	Title       string
	Description string
	DueDate     *time.Time
	Priority    string
	AssigneeID  *uint
}

type UpdateTaskInput struct {
	Title       *string
	Description *string
	DueDate     *time.Time
	Priority    *string
	AssigneeID  *uint
}

type UpdateTaskStatusInput struct {
	Status string
}

type TaskService struct {
	tasks    repository.TaskRepository
	projects repository.ProjectRepository
}

func NewTaskService(tasks repository.TaskRepository, projects repository.ProjectRepository) *TaskService {
	return &TaskService{
		tasks:    tasks,
		projects: projects,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, userID uint, input CreateTaskInput) (*models.Task, error) {
	if err := s.requireAdmin(ctx, input.ProjectID, userID); err != nil {
		return nil, err
	}

	if input.AssigneeID != nil {
		if _, err := s.projects.GetMembership(ctx, input.ProjectID, *input.AssigneeID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperror.BadRequest("assignee must be a member of the project")
			}
			return nil, apperror.Internal("failed to verify assignee membership")
		}
	}

	task := &models.Task{
		ProjectID:   input.ProjectID,
		CreatorID:   userID,
		AssigneeID:  input.AssigneeID,
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		DueDate:     input.DueDate,
		Priority:    models.TaskPriority(input.Priority),
		Status:      models.TaskStatusTodo,
	}

	if err := s.tasks.Create(ctx, task); err != nil {
		return nil, apperror.Internal("failed to create task")
	}

	createdTask, err := s.tasks.GetByID(ctx, task.ID)
	if err != nil {
		return nil, apperror.Internal("failed to reload task")
	}

	return createdTask, nil
}

func (s *TaskService) ListTasks(ctx context.Context, userID, projectID uint) ([]models.Task, error) {
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

	tasks, err := s.tasks.ListByProject(ctx, projectID, assigneeFilter)
	if err != nil {
		return nil, apperror.Internal("failed to list tasks")
	}

	return tasks, nil
}

func (s *TaskService) GetTask(ctx context.Context, userID, taskID uint) (*models.Task, error) {
	task, err := s.tasks.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("task not found")
		}
		return nil, apperror.Internal("failed to fetch task")
	}

	if err := s.ensureTaskAccess(ctx, task, userID, true); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, userID, taskID uint, input UpdateTaskInput) (*models.Task, error) {
	task, err := s.tasks.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("task not found")
		}
		return nil, apperror.Internal("failed to fetch task")
	}

	if err := s.requireAdmin(ctx, task.ProjectID, userID); err != nil {
		return nil, err
	}

	if input.Title != nil {
		task.Title = strings.TrimSpace(*input.Title)
	}
	if input.Description != nil {
		task.Description = strings.TrimSpace(*input.Description)
	}
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}
	if input.Priority != nil {
		task.Priority = models.TaskPriority(*input.Priority)
	}
	if input.AssigneeID != nil {
		if *input.AssigneeID == 0 {
			task.AssigneeID = nil
		} else if _, err := s.projects.GetMembership(ctx, task.ProjectID, *input.AssigneeID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperror.BadRequest("assignee must be a member of the project")
			}
			return nil, apperror.Internal("failed to verify assignee membership")
		} else {
			task.AssigneeID = input.AssigneeID
		}
	}

	if err := s.tasks.Update(ctx, task); err != nil {
		return nil, apperror.Internal("failed to update task")
	}

	updatedTask, err := s.tasks.GetByID(ctx, task.ID)
	if err != nil {
		return nil, apperror.Internal("failed to reload task")
	}

	return updatedTask, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, userID, taskID uint) error {
	task, err := s.tasks.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("task not found")
		}
		return apperror.Internal("failed to fetch task")
	}

	if err := s.requireAdmin(ctx, task.ProjectID, userID); err != nil {
		return err
	}

	if err := s.tasks.Delete(ctx, taskID); err != nil {
		return apperror.Internal("failed to delete task")
	}

	return nil
}

func (s *TaskService) UpdateTaskStatus(ctx context.Context, userID, taskID uint, input UpdateTaskStatusInput) (*models.Task, error) {
	task, err := s.tasks.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("task not found")
		}
		return nil, apperror.Internal("failed to fetch task")
	}

	if err := s.ensureTaskAccess(ctx, task, userID, true); err != nil {
		return nil, err
	}

	task.Status = models.TaskStatus(input.Status)
	if err := s.tasks.Update(ctx, task); err != nil {
		return nil, apperror.Internal("failed to update task status")
	}

	updatedTask, err := s.tasks.GetByID(ctx, task.ID)
	if err != nil {
		return nil, apperror.Internal("failed to reload task")
	}

	return updatedTask, nil
}

func (s *TaskService) ensureTaskAccess(ctx context.Context, task *models.Task, userID uint, allowAssignedMember bool) error {
	member, err := s.projects.GetMembership(ctx, task.ProjectID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.Forbidden("you do not belong to this project")
		}
		return apperror.Internal("failed to verify project membership")
	}

	if member.Role == models.ProjectRoleAdmin {
		return nil
	}

	if allowAssignedMember && task.AssigneeID != nil && *task.AssigneeID == userID {
		return nil
	}

	return apperror.Forbidden("you can only access tasks assigned to you")
}

func (s *TaskService) requireAdmin(ctx context.Context, projectID, userID uint) error {
	member, err := s.projects.GetMembership(ctx, projectID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.Forbidden("you do not belong to this project")
		}
		return apperror.Internal("failed to verify project membership")
	}

	if member.Role != models.ProjectRoleAdmin {
		return apperror.Forbidden("admin access is required for this project")
	}

	return nil
}
