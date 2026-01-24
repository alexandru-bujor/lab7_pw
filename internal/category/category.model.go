package category

import "time"

type Category struct {
	ID        int       `gorm:"column=id;primaryKey" json:"id"`
	Name      string    `gorm:"column=name" json:"name"`
	Slug      string    `gorm:"column=slug" json:"slug"`
	ParentID  *int      `gorm:"column=parent_id" json:"parent_id,omitempty"`
	CreatedAt time.Time `gorm:"column=created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column=updated_at" json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}
