package make

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, m)
	case string:
		return json.Unmarshal([]byte(v), m)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

type Make struct {
	ID   uint            `gorm:"primaryKey" json:"id"`
	Name MultilingualName `gorm:"type:json" json:"name"`
	Models []Model       `gorm:"foreignKey:MakeID" json:"models"`
}

type Model struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	MakeID uint   `json:"make_id"`
	Name   string `json:"name"`
}

func (Make) TableName() string {
	return "makes"
}

func (Model) TableName() string {
	return "models"
}
