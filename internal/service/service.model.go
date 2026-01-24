package service

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// PricesMap maps repair_type_id (string) to price (int)
type PricesMap map[string]int

// Value implements driver.Valuer for GORM/SQL
func (p PricesMap) Value() (driver.Value, error) {
	if p == nil {
		return "{}", nil
	}
	return json.Marshal(p)
}

// Scan implements sql.Scanner for GORM/SQL
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
	ID          uint      `gorm:"primaryKey" json:"id"`
	DeviceName  string    `json:"device_name"`
	Display     int       `json:"display"`      // Legacy field, kept for backward compatibility
	BackGlass   int       `json:"back_glass"`    // Legacy field, kept for backward compatibility
	Battery     int       `json:"battery"`       // Legacy field, kept for backward compatibility
	Prices      PricesMap `gorm:"type:json" json:"prices"` // Maps repair_type_id -> price
	RepairTypes string    `json:"repair_types"`  // JSON array stored as string
	ServicePart int       `json:"service_part"`  // 1 or 2
	PhotoURL    string    `json:"photo_url"`
	IsVisibleInSecondary bool `gorm:"column:is_visible_in_secondary;default:false" json:"is_visible_in_secondary"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func (Service) TableName() string {
	return "services"
}
