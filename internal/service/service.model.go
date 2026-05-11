package service

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type PricesMap map[string]int

func (p PricesMap) Value() (driver.Value, error) {
	if p == nil {
		return "{}", nil
	}
	return json.Marshal(p)
}

func (p *PricesMap) Scan(value interface{}) error {
	if value == nil {
		*p = make(PricesMap)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, p)
	case string:
		return json.Unmarshal([]byte(v), p)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

type Service struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	DeviceName           string    `json:"device_name"`
	Display              int       `json:"display"`
	BackGlass            int       `json:"back_glass"`
	Battery              int       `json:"battery"`
	Prices               PricesMap `gorm:"type:json" json:"prices"`
	RepairTypes          string    `json:"repair_types"`
	ServicePart          int       `json:"service_part"`
	PhotoURL             string    `json:"photo_url"`
	IsVisibleInSecondary bool      `gorm:"column:is_visible_in_secondary;default:false" json:"is_visible_in_secondary"`
	CreatedAt            string    `json:"created_at"`
	UpdatedAt            string    `json:"updated_at"`
}

func (Service) TableName() string {
	return "services"
}
