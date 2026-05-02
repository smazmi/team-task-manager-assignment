package service

import (
	"context"
	"errors"
	"strings"

	"github.com/smazmi/team-task-manager-assignment/backend/internal/models"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/repository"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/apperror"
	"gorm.io/gorm"
)

type CreateProjectInput struct {
	Name        string
	Description string
}

type UpdateProjectInput struct {
	Name        *string
	Description *string
}

type AddProjectMemberInput struct {
	Email string
	Role  models.ProjectRole
}

type ProjectService struct {
	projects repository.ProjectRepository
	users    repository.UserRepository
}

func NewProjectService(projects repository.ProjectRepository, users repository.UserRepository) *ProjectService {
	return &ProjectService{
		projects: projects,
		users:    users,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, creatorID uint, input CreateProjectInput) (*models.Project, error) {
	project := &models.Project{
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		CreatorID:   creatorID,
		Members: []models.ProjectMember{
			{
				UserID: creatorID,
				Role:   models.ProjectRoleAdmin,
			},
		},
	}

	if err := s.projects.Create(ctx, project); err != nil {
		return nil, apperror.Internal("failed to create project")
	}

	createdProject, err := s.projects.GetByID(ctx, project.ID)
	if err != nil {
		return nil, apperror.Internal("failed to reload project")
	}

	return createdProject, nil
}

func (s *ProjectService) ListProjects(ctx context.Context, userID uint) ([]models.Project, error) {
	projects, err := s.projects.ListByUserID(ctx, userID)
	if err != nil {
		return nil, apperror.Internal("failed to list projects")
	}

	return projects, nil
}

func (s *ProjectService) GetProject(ctx context.Context, userID, projectID uint) (*models.Project, error) {
	if _, err := s.requireMembership(ctx, projectID, userID); err != nil {
		return nil, err
	}

	project, err := s.projects.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("project not found")
		}
		return nil, apperror.Internal("failed to fetch project")
	}

	return project, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, userID, projectID uint, input UpdateProjectInput) (*models.Project, error) {
	if err := s.requireAdmin(ctx, projectID, userID); err != nil {
		return nil, err
	}

	project, err := s.projects.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("project not found")
		}
		return nil, apperror.Internal("failed to fetch project")
	}

	if input.Name != nil {
		project.Name = strings.TrimSpace(*input.Name)
	}
	if input.Description != nil {
		project.Description = strings.TrimSpace(*input.Description)
	}

	if err := s.projects.Update(ctx, project); err != nil {
		return nil, apperror.Internal("failed to update project")
	}

	updatedProject, err := s.projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, apperror.Internal("failed to reload project")
	}

	return updatedProject, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, userID, projectID uint) error {
	if err := s.requireAdmin(ctx, projectID, userID); err != nil {
		return err
	}

	if err := s.projects.Delete(ctx, projectID); err != nil {
		return apperror.Internal("failed to delete project")
	}

	return nil
}

func (s *ProjectService) ListMembers(ctx context.Context, userID, projectID uint) ([]models.ProjectMember, error) {
	if _, err := s.requireMembership(ctx, projectID, userID); err != nil {
		return nil, err
	}

	members, err := s.projects.ListMembers(ctx, projectID)
	if err != nil {
		return nil, apperror.Internal("failed to list project members")
	}

	return members, nil
}

func (s *ProjectService) AddMember(ctx context.Context, actingUserID, projectID uint, input AddProjectMemberInput) (*models.Project, error) {
	if err := s.requireAdmin(ctx, projectID, actingUserID); err != nil {
		return nil, err
	}

	if _, err := s.projects.GetByID(ctx, projectID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("project not found")
		}
		return nil, apperror.Internal("failed to fetch project")
	}

	user, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("user not found with that email")
		}
		return nil, apperror.Internal("failed to fetch user")
	}

	existingMember, err := s.projects.GetMembership(ctx, projectID, user.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.Internal("failed to verify project membership")
	}

	if existingMember != nil && existingMember.Role == models.ProjectRoleAdmin && input.Role != models.ProjectRoleAdmin {
		adminCount, err := s.projects.CountAdmins(ctx, projectID)
		if err != nil {
			return nil, apperror.Internal("failed to verify admin count")
		}
		if adminCount == 1 {
			return nil, apperror.BadRequest("project must have at least one admin")
		}
	}

	member := &models.ProjectMember{
		ProjectID: projectID,
		UserID:    user.ID,
		Role:      input.Role,
	}

	if err := s.projects.UpsertMember(ctx, member); err != nil {
		return nil, apperror.Internal("failed to add project member")
	}

	project, err := s.projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, apperror.Internal("failed to reload project")
	}

	return project, nil
}

func (s *ProjectService) RemoveMember(ctx context.Context, actingUserID, projectID, targetUserID uint) error {
	if err := s.requireAdmin(ctx, projectID, actingUserID); err != nil {
		return err
	}

	member, err := s.projects.GetMembership(ctx, projectID, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("project member not found")
		}
		return apperror.Internal("failed to fetch project member")
	}

	if member.Role == models.ProjectRoleAdmin {
		adminCount, err := s.projects.CountAdmins(ctx, projectID)
		if err != nil {
			return apperror.Internal("failed to verify admin count")
		}
		if adminCount == 1 {
			return apperror.BadRequest("project must have at least one admin")
		}
	}

	if err := s.projects.RemoveMember(ctx, projectID, targetUserID); err != nil {
		return apperror.Internal("failed to remove project member")
	}

	return nil
}

func (s *ProjectService) requireMembership(ctx context.Context, projectID, userID uint) (*models.ProjectMember, error) {
	member, err := s.projects.GetMembership(ctx, projectID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Forbidden("you do not belong to this project")
		}
		return nil, apperror.Internal("failed to verify project membership")
	}

	return member, nil
}

func (s *ProjectService) requireAdmin(ctx context.Context, projectID, userID uint) error {
	member, err := s.requireMembership(ctx, projectID, userID)
	if err != nil {
		return err
	}

	if member.Role != models.ProjectRoleAdmin {
		return apperror.Forbidden("admin access is required for this project")
	}

	return nil
}
