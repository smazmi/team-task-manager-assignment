package dto

import (
	"time"

	"github.com/smazmi/team-task-manager-assignment/backend/internal/models"
)

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=120"`
	Description string `json:"description" binding:"omitempty,max=2000"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=3,max=120"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
}

type AddProjectMemberRequest struct {
	UserID uint   `json:"user_id" binding:"required,gt=0"`
	Role   string `json:"role" binding:"required,oneof=admin member"`
}

type ProjectMemberResponse struct {
	ID        uint         `json:"id"`
	ProjectID uint         `json:"project_id"`
	UserID    uint         `json:"user_id"`
	Role      string       `json:"role"`
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type ProjectResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	CreatorID   uint                    `json:"creator_id"`
	Creator     UserResponse            `json:"creator"`
	Members     []ProjectMemberResponse `json:"members"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

func NewProjectMemberResponse(member models.ProjectMember) ProjectMemberResponse {
	return ProjectMemberResponse{
		ID:        member.ID,
		ProjectID: member.ProjectID,
		UserID:    member.UserID,
		Role:      string(member.Role),
		User:      NewUserResponse(&member.User),
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}

func NewProjectResponse(project *models.Project) ProjectResponse {
	members := make([]ProjectMemberResponse, 0, len(project.Members))
	for _, member := range project.Members {
		members = append(members, NewProjectMemberResponse(member))
	}

	return ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		CreatorID:   project.CreatorID,
		Creator:     NewUserResponse(&project.Creator),
		Members:     members,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
}

func NewProjectResponses(projects []models.Project) []ProjectResponse {
	responses := make([]ProjectResponse, 0, len(projects))
	for _, project := range projects {
		projectCopy := project
		responses = append(responses, NewProjectResponse(&projectCopy))
	}

	return responses
}

func NewProjectMemberResponses(members []models.ProjectMember) []ProjectMemberResponse {
	responses := make([]ProjectMemberResponse, 0, len(members))
	for _, member := range members {
		responses = append(responses, NewProjectMemberResponse(member))
	}

	return responses
}
