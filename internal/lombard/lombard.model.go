package lombard

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// LombardProduct represents a product/item in a lombard request
type LombardProduct struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Condition   string `json:"condition"`
}

// LombardProducts is a slice of LombardProduct for JSON storage
type LombardProducts []LombardProduct

// Value implements driver.Valuer for GORM/SQL
func (p LombardProducts) Value() (driver.Value, error) {
	if len(p) == 0 {
		return "[]", nil
	}
	return json.Marshal(p)
}

// Scan implements sql.Scanner for GORM/SQL
func (p *LombardProducts) Scan(value interface{}) error {
	if value == nil {
		*p = LombardProducts{}
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

// LombardRequest represents a lombard request from a customer
type LombardRequest struct {
	ID           int             `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CustomerName string          `gorm:"column:customer_name;not null" json:"customer_name"`
	CustomerPhone string         `gorm:"column:customer_phone;not null" json:"customer_phone"`
	RequestType  string          `gorm:"column:request_type" json:"request_type"` // 'electronics' or 'gold'
	Products     LombardProducts `gorm:"column:products;type:json" json:"products"`
	Status       string          `gorm:"column:status;default:'new'" json:"status"` // new | old | seen
	IsVisibleInSecondary bool    `gorm:"column:is_visible_in_secondary;default:false" json:"is_visible_in_secondary"`
	CreatedAt    time.Time       `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (LombardRequest) TableName() string {
	return "lombard_requests"
}

// CreateLombardRequestInput for creating lombard requests from website
type CreateLombardRequestInput struct {
	CustomerName  string          `json:"customer_name" binding:"required"`
	CustomerPhone string          `json:"customer_phone" binding:"required"`
	RequestType   string          `json:"request_type"` // 'electronics' or 'gold'
	Products      LombardProducts `json:"products" binding:"required"`
}

// UpdateStatusInput for updating request status
type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
}

