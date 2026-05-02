package models

import "time"

type User struct {
	ID                 uint            `gorm:"primaryKey" json:"id"`
	Name               string          `gorm:"size:120;not null" json:"name"`
	Email              string          `gorm:"size:255;not null;uniqueIndex" json:"email"`
	PasswordHash       string          `gorm:"column:password_hash;size:255;not null" json:"-"`
	CreatedProjects    []Project       `gorm:"foreignKey:CreatorID" json:"-"`
	ProjectMemberships []ProjectMember `gorm:"foreignKey:UserID" json:"-"`
	AssignedTasks      []Task          `gorm:"foreignKey:AssigneeID" json:"-"`
	CreatedTasks       []Task          `gorm:"foreignKey:CreatorID" json:"-"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}
