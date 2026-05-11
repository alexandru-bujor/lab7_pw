package category

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type MultilingualName struct {
	RO string `json:"ro"`
	RU string `json:"ru"`
	EN string `json:"en"`
}

func (m MultilingualName) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MultilingualName) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}

	if len(bytes) == 0 {
		return nil
	}

	// Try to unmarshal as JSON
	if err := json.Unmarshal(bytes, m); err != nil {
		// Fallback: if not valid JSON, treat as a legacy plain string for all languages
		plainStr := string(bytes)
		m.RO = plainStr
		m.EN = plainStr
		m.RU = plainStr
	}
	return nil
}

type Category struct {
	ID        int              `gorm:"column:id" json:"id"`
	Name      MultilingualName `gorm:"column:name;type:json" json:"name"`
	Slug      string           `gorm:"column:slug" json:"slug"`
	ParentID  *int             `gorm:"column:parent_id" json:"parent_id,omitempty"`
	CreatedAt time.Time        `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time        `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}
