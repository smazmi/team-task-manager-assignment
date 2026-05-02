package models

import "time"

type Project struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"size:120;not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	CreatorID   uint            `gorm:"column:creator_id;not null;index" json:"creator_id"`
	Creator     User            `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"creator"`
	Members     []ProjectMember `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"members"`
	Tasks       []Task          `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tasks"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
