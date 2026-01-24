package inventory

import (
	"time"
)

type InventoryItem struct {
	ID            int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductID     int        `gorm:"column:product_id;not null" json:"product_id"`
	Section       string     `gorm:"column:section;not null" json:"section"` // e.g. "Apple SH", "Samsung", "Others", "PC"
	VariantColor  string     `gorm:"column:variant_color" json:"variant_color"`
	VariantStorage string    `gorm:"column:variant_storage" json:"variant_storage"`
	VariantRAM    string     `gorm:"column:variant_ram" json:"variant_ram"`
	Serial        string     `gorm:"column:serial" json:"serial"`
	BuyPrice      float64    `gorm:"column:buy_price" json:"buy_price"`
	SellPrice     float64    `gorm:"column:sell_price" json:"sell_price"`
	Code          string     `gorm:"column:code" json:"code"`
	Source        string     `gorm:"column:source" json:"source"`
	Status        string     `gorm:"column:status;default:'in_stock'" json:"status"` // in_stock | reserved | sold | maintenance
	ClientID      *int       `gorm:"column:client_id" json:"client_id,omitempty"`
	AssignedAt    *time.Time `gorm:"column:assigned_at" json:"assigned_at,omitempty"`
	SoldAt        *time.Time `gorm:"column:sold_at" json:"sold_at,omitempty"`
	Notes         string     `gorm:"column:notes" json:"notes"`
	IsVisibleInSecondary bool `gorm:"column:is_visible_in_secondary;default:false" json:"is_visible_in_secondary"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (InventoryItem) TableName() string { return "inventory_items" }
