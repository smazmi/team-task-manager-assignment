package dto

import (
	"time"

	"github.com/smazmi/team-task-manager-assignment/backend/internal/models"
)

type CreateTaskRequest struct {
	ProjectID   uint       `json:"project_id" binding:"required,gt=0"`
	Title       string     `json:"title" binding:"required,min=1,max=200"`
	Description string     `json:"description" binding:"omitempty,max=4000"`
	DueDate     *time.Time `json:"due_date"`
	Priority    string     `json:"priority" binding:"required,oneof=low medium high"`
	AssigneeID  *uint      `json:"assignee_id" binding:"omitempty,gt=0"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=200"`
	Description *string    `json:"description" binding:"omitempty,max=4000"`
	DueDate     *time.Time `json:"due_date"`
	Priority    *string    `json:"priority" binding:"omitempty,oneof=low medium high"`
	AssigneeID  *uint      `json:"assignee_id" binding:"omitempty,gte=0"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=todo in_progress done"`
}

type ListTasksQuery struct {
	ProjectID uint `form:"project_id" binding:"required,gt=0"`
}

type TaskResponse struct {
	ID          uint          `json:"id"`
	ProjectID   uint          `json:"project_id"`
	CreatorID   uint          `json:"creator_id"`
	AssigneeID  *uint         `json:"assignee_id,omitempty"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	DueDate     *time.Time    `json:"due_date,omitempty"`
	Priority    string        `json:"priority"`
	Status      string        `json:"status"`
	Creator     UserResponse  `json:"creator"`
	Assignee    *UserResponse `json:"assignee,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func NewTaskResponse(task *models.Task) TaskResponse {
	var assignee *UserResponse
	if task.Assignee != nil {
		assigneeResponse := NewUserResponse(task.Assignee)
		assignee = &assigneeResponse
	}

	return TaskResponse{
		ID:          task.ID,
		ProjectID:   task.ProjectID,
		CreatorID:   task.CreatorID,
		AssigneeID:  task.AssigneeID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Priority:    string(task.Priority),
		Status:      string(task.Status),
		Creator:     NewUserResponse(&task.Creator),
		Assignee:    assignee,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func NewTaskResponses(tasks []models.Task) []TaskResponse {
	responses := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		taskCopy := task
		responses = append(responses, NewTaskResponse(&taskCopy))
	}

	return responses
}
