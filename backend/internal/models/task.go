package models

import "time"

type TaskPriority string
type TaskStatus string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

type Task struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	ProjectID   uint         `gorm:"column:project_id;not null;index" json:"project_id"`
	CreatorID   uint         `gorm:"column:creator_id;not null;index" json:"creator_id"`
	AssigneeID  *uint        `gorm:"column:assignee_id;index" json:"assignee_id,omitempty"`
	Title       string       `gorm:"size:200;not null" json:"title"`
	Description string       `gorm:"type:text" json:"description"`
	DueDate     *time.Time   `gorm:"column:due_date" json:"due_date,omitempty"`
	Priority    TaskPriority `gorm:"type:varchar(20);not null;default:'medium';check:priority IN ('low','medium','high')" json:"priority"`
	Status      TaskStatus   `gorm:"type:varchar(20);not null;default:'todo';check:status IN ('todo','in_progress','done')" json:"status"`
	Project     Project      `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Creator     User         `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"creator"`
	Assignee    *User        `gorm:"foreignKey:AssigneeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"assignee,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
