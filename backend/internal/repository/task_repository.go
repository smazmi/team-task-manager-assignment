package repository

import (
	"context"
	"time"

	"github.com/smazmi/team-task-manager-assignment/backend/internal/models"
	"gorm.io/gorm"
)

type TaskCountByUser struct {
	UserID   uint
	UserName string
	Count    int64
}

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, taskID uint) (*models.Task, error)
	ListByProject(ctx context.Context, projectID uint, assigneeID *uint) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, taskID uint) error
	CountByProjectAndAssignee(ctx context.Context, projectID uint, assigneeID *uint) (int64, error)
	GetStatusCounts(ctx context.Context, projectID uint, assigneeID *uint) (map[models.TaskStatus]int64, error)
	GetTasksPerUser(ctx context.Context, projectID uint, assigneeID *uint) ([]TaskCountByUser, error)
	CountOverdue(ctx context.Context, projectID uint, assigneeID *uint, currentTime time.Time) (int64, error)
}

type GormTaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &GormTaskRepository{db: db}
}

func (r *GormTaskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *GormTaskRepository) GetByID(ctx context.Context, taskID uint) (*models.Task, error) {
	var task models.Task
	if err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Assignee").
		First(&task, taskID).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *GormTaskRepository) ListByProject(ctx context.Context, projectID uint, assigneeID *uint) ([]models.Task, error) {
	var tasks []models.Task
	query := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Preload("Creator").
		Preload("Assignee")

	if assigneeID != nil {
		query = query.Where("assignee_id = ?", *assigneeID)
	}

	if err := query.
		Order("due_date IS NULL, due_date ASC").
		Order("created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *GormTaskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *GormTaskRepository) Delete(ctx context.Context, taskID uint) error {
	return r.db.WithContext(ctx).Delete(&models.Task{}, taskID).Error
}

func (r *GormTaskRepository) CountByProjectAndAssignee(ctx context.Context, projectID uint, assigneeID *uint) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Task{}).Where("project_id = ?", projectID)
	if assigneeID != nil {
		query = query.Where("assignee_id = ?", *assigneeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *GormTaskRepository) GetStatusCounts(ctx context.Context, projectID uint, assigneeID *uint) (map[models.TaskStatus]int64, error) {
	type statusCountRow struct {
		Status models.TaskStatus
		Count  int64
	}

	query := r.db.WithContext(ctx).
		Model(&models.Task{}).
		Select("status, COUNT(*) AS count").
		Where("project_id = ?", projectID)

	if assigneeID != nil {
		query = query.Where("assignee_id = ?", *assigneeID)
	}

	var rows []statusCountRow
	if err := query.Group("status").Scan(&rows).Error; err != nil {
		return nil, err
	}

	counts := make(map[models.TaskStatus]int64, len(rows))
	for _, row := range rows {
		counts[row.Status] = row.Count
	}

	return counts, nil
}

func (r *GormTaskRepository) GetTasksPerUser(ctx context.Context, projectID uint, assigneeID *uint) ([]TaskCountByUser, error) {
	var rows []TaskCountByUser
	query := r.db.WithContext(ctx).
		Table("tasks").
		Select("users.id AS user_id, users.name AS user_name, COUNT(tasks.id) AS count").
		Joins("JOIN users ON users.id = tasks.assignee_id").
		Where("tasks.project_id = ?", projectID)

	if assigneeID != nil {
		query = query.Where("tasks.assignee_id = ?", *assigneeID)
	}

	if err := query.
		Group("users.id, users.name").
		Order("count DESC, users.name ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *GormTaskRepository) CountOverdue(ctx context.Context, projectID uint, assigneeID *uint, currentTime time.Time) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("project_id = ?", projectID).
		Where("due_date IS NOT NULL").
		Where("due_date < ?", currentTime).
		Where("status <> ?", models.TaskStatusDone)

	if assigneeID != nil {
		query = query.Where("assignee_id = ?", *assigneeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
