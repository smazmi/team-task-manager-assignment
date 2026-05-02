package models

import "time"

type ProjectRole string

const (
	ProjectRoleAdmin  ProjectRole = "admin"
	ProjectRoleMember ProjectRole = "member"
)

type ProjectMember struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	ProjectID uint        `gorm:"column:project_id;not null;index;uniqueIndex:idx_project_member" json:"project_id"`
	UserID    uint        `gorm:"column:user_id;not null;index;uniqueIndex:idx_project_member" json:"user_id"`
	Role      ProjectRole `gorm:"type:varchar(20);not null;default:'member';check:role IN ('admin','member')" json:"role"`
	Project   Project     `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	User      User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
