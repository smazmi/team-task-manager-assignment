package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/dto"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/service"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/apperror"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/response"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var request dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	task, err := h.taskService.CreateTask(c.Request.Context(), userID, service.CreateTaskInput{
		ProjectID:   request.ProjectID,
		Title:       request.Title,
		Description: request.Description,
		DueDate:     request.DueDate,
		Priority:    request.Priority,
		AssigneeID:  request.AssigneeID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "task created", dto.NewTaskResponse(task))
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	var query dto.ListTasksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	tasks, err := h.taskService.ListTasks(c.Request.Context(), userID, query.ProjectID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "tasks fetched", dto.NewTaskResponses(tasks))
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	taskID, err := parseUintParam(c, "taskId")
	if err != nil {
		response.Error(c, err)
		return
	}

	task, err := h.taskService.GetTask(c.Request.Context(), userID, taskID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "task fetched", dto.NewTaskResponse(task))
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var request dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	taskID, err := parseUintParam(c, "taskId")
	if err != nil {
		response.Error(c, err)
		return
	}

	task, err := h.taskService.UpdateTask(c.Request.Context(), userID, taskID, service.UpdateTaskInput{
		Title:       request.Title,
		Description: request.Description,
		DueDate:     request.DueDate,
		Priority:    request.Priority,
		AssigneeID:  request.AssigneeID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "task updated", dto.NewTaskResponse(task))
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	taskID, err := parseUintParam(c, "taskId")
	if err != nil {
		response.Error(c, err)
		return
	}

	if err := h.taskService.DeleteTask(c.Request.Context(), userID, taskID); err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c, "task deleted")
}

func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	var request dto.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	taskID, err := parseUintParam(c, "taskId")
	if err != nil {
		response.Error(c, err)
		return
	}

	task, err := h.taskService.UpdateTaskStatus(c.Request.Context(), userID, taskID, service.UpdateTaskStatusInput{
		Status: request.Status,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "task status updated", dto.NewTaskResponse(task))
}
