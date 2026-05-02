package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/dto"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/models"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/service"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/apperror"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/response"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var request dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	project, err := h.projectService.CreateProject(c.Request.Context(), userID, service.CreateProjectInput{
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "project created", dto.NewProjectResponse(project))
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	projects, err := h.projectService.ListProjects(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "projects fetched", dto.NewProjectResponses(projects))
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	projectID, err := parseUintParam(c, "projectId")
	if err != nil {
		response.Error(c, err)
		return
	}

	project, err := h.projectService.GetProject(c.Request.Context(), userID, projectID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "project fetched", dto.NewProjectResponse(project))
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	var request dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	projectID, err := parseUintParam(c, "projectId")
	if err != nil {
		response.Error(c, err)
		return
	}

	project, err := h.projectService.UpdateProject(c.Request.Context(), userID, projectID, service.UpdateProjectInput{
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "project updated", dto.NewProjectResponse(project))
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	projectID, err := parseUintParam(c, "projectId")
	if err != nil {
		response.Error(c, err)
		return
	}

	if err := h.projectService.DeleteProject(c.Request.Context(), userID, projectID); err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c, "project deleted")
}

func (h *ProjectHandler) ListMembers(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	projectID, err := parseUintParam(c, "projectId")
	if err != nil {
		response.Error(c, err)
		return
	}

	members, err := h.projectService.ListMembers(c.Request.Context(), userID, projectID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "project members fetched", dto.NewProjectMemberResponses(members))
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	var request dto.AddProjectMemberRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	projectID, err := parseUintParam(c, "projectId")
	if err != nil {
		response.Error(c, err)
		return
	}

	project, err := h.projectService.AddMember(c.Request.Context(), userID, projectID, service.AddProjectMemberInput{
		Email: request.Email,
		Role:  models.ProjectRole(request.Role),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "project member added", dto.NewProjectResponse(project))
}

func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	userID, err := currentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	projectID, err := parseUintParam(c, "projectId")
	if err != nil {
		response.Error(c, err)
		return
	}

	targetUserID, err := parseUintParam(c, "userId")
	if err != nil {
		response.Error(c, err)
		return
	}

	if err := h.projectService.RemoveMember(c.Request.Context(), userID, projectID, targetUserID); err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c, "project member removed")
}
