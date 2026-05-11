package category

import (
	"time"

	"gorm.io/datatypes"
)

type CategoryTemplate struct {
	ID         int               `gorm:"column=id;primaryKey" json:"id"`
	CategoryID int               `gorm:"column=category_id;index;not null" json:"category_id"`
	Name       string            `gorm:"column=name" json:"name"`
	Specs      datatypes.JSONMap `gorm:"column=specs;type:json" json:"specs"`
	CreatedAt  time.Time         `gorm:"column=created_at" json:"created_at"`
	UpdatedAt  time.Time         `gorm:"column=updated_at" json:"updated_at"`
}

func (CategoryTemplate) TableName() string {
	return "category_templates"
}
