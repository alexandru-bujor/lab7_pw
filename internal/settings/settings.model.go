package settings

import (
	"time"
)

type Setting struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Key       string    `gorm:"column:key;unique;not null" json:"key"`
	Value     string    `gorm:"column:value" json:"value"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Setting) TableName() string {
	return "settings"
}

const (
	SettingKeySecondaryPanelPassword = "secondary_panel_password"
)
