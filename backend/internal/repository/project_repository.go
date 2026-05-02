package repository

import (
	"context"
	"time"

	"github.com/smazmi/team-task-manager-assignment/backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, projectID uint) (*models.Project, error)
	ListByUserID(ctx context.Context, userID uint) ([]models.Project, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, projectID uint) error
	GetMembership(ctx context.Context, projectID, userID uint) (*models.ProjectMember, error)
	ListMembers(ctx context.Context, projectID uint) ([]models.ProjectMember, error)
	UpsertMember(ctx context.Context, member *models.ProjectMember) error
	RemoveMember(ctx context.Context, projectID, userID uint) error
	UserHasRole(ctx context.Context, projectID, userID uint, roles ...models.ProjectRole) (bool, error)
	CountAdmins(ctx context.Context, projectID uint) (int64, error)
}

type GormProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &GormProjectRepository{db: db}
}

func (r *GormProjectRepository) Create(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *GormProjectRepository) GetByID(ctx context.Context, projectID uint) (*models.Project, error) {
	var project models.Project
	if err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Members.User").
		First(&project, projectID).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *GormProjectRepository) ListByUserID(ctx context.Context, userID uint) ([]models.Project, error) {
	var projects []models.Project
	if err := r.db.WithContext(ctx).
		Model(&models.Project{}).
		Joins("JOIN project_members ON project_members.project_id = projects.id").
		Where("project_members.user_id = ?", userID).
		Preload("Creator").
		Preload("Members.User").
		Order("projects.created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *GormProjectRepository) Update(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *GormProjectRepository) Delete(ctx context.Context, projectID uint) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, projectID).Error
}

func (r *GormProjectRepository) GetMembership(ctx context.Context, projectID, userID uint) (*models.ProjectMember, error) {
	var member models.ProjectMember
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("project_id = ? AND user_id = ?", projectID, userID).
		First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

func (r *GormProjectRepository) ListMembers(ctx context.Context, projectID uint) ([]models.ProjectMember, error) {
	var members []models.ProjectMember
	if err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Preload("User").
		Order("created_at ASC").
		Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

func (r *GormProjectRepository) UpsertMember(ctx context.Context, member *models.ProjectMember) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "project_id"},
				{Name: "user_id"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"role":       member.Role,
				"updated_at": time.Now(),
			}),
		}).
		Create(member).Error
}

func (r *GormProjectRepository) RemoveMember(ctx context.Context, projectID, userID uint) error {
	return r.db.WithContext(ctx).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&models.ProjectMember{}).Error
}

func (r *GormProjectRepository) UserHasRole(ctx context.Context, projectID, userID uint, roles ...models.ProjectRole) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID)

	if len(roles) > 0 {
		query = query.Where("role IN ?", roles)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *GormProjectRepository) CountAdmins(ctx context.Context, projectID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.ProjectMember{}).
		Where("project_id = ? AND role = ?", projectID, models.ProjectRoleAdmin).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
