package order

import (
	"time"
)

type Order struct {
	ID                   int        `gorm:"column=id;primaryKey;autoIncrement" json:"id"`
	ClientID             *int       `gorm:"column=client_id" json:"client_id"`
	ClientName           string     `gorm:"column=client_name" json:"client_name"`
	ClientEmail          string     `gorm:"column=client_email" json:"client_email"`
	ClientPhone          string     `gorm:"column=client_phone" json:"client_phone"`
	InventoryItemID      *int       `gorm:"column=inventory_item_id" json:"inventory_item_id,omitempty"`
	ProductID            int        `gorm:"column=product_id" json:"product_id"`
	ProductTitle         string     `gorm:"column=product_title" json:"product_title"`
	VariantColor         string     `gorm:"column=variant_color" json:"variant_color"`
	VariantStorage       string     `gorm:"column=variant_storage" json:"variant_storage"`
	VariantRAM           string     `gorm:"column=variant_ram" json:"variant_ram"`
	Quantity             int        `gorm:"column=quantity;default:1" json:"quantity"`
	SellPrice            float64    `gorm:"column=sell_price" json:"sell_price"`
	BuyPrice             float64    `gorm:"column=buy_price;default:0" json:"buy_price"`
	Profit               float64    `gorm:"column=profit;default:0" json:"profit"`
	Status               string     `gorm:"column=status;default:pending" json:"status"`
	OrderSource          string     `gorm:"column=order_source;default:website" json:"order_source"`
	AccountingType       string     `gorm:"column=accounting_type;default:on_book" json:"accounting_type"`
	IsVisibleInSecondary bool       `gorm:"column=is_visible_in_secondary;default:false" json:"is_visible_in_secondary"`
	ManagerID            *int       `gorm:"column=manager_id" json:"manager_id"`
	ManagerName          string     `gorm:"->" json:"manager_name,omitempty"`
	ProductBrandName     string     `gorm:"->" json:"product_brand_name,omitempty"`
	ProductModel         string     `gorm:"->" json:"product_model,omitempty"`
	SoldAt               *time.Time `gorm:"column=sold_at" json:"sold_at,omitempty"`
	CreatedAt            time.Time  `gorm:"column=created_at" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"column=updated_at" json:"updated_at"`
}

func (Order) TableName() string { return "orders" }

type CreateOrderInput struct {
	ClientName     string  `json:"client_name" binding:"required"`
	ClientEmail    string  `json:"client_email" binding:"required"`
	ClientPhone    string  `json:"client_phone" binding:"required"`
	ProductID      int     `json:"product_id" binding:"required,min=1"`
	ProductTitle   string  `json:"product_title" binding:"required"`
	VariantColor   string  `json:"variant_color"`
	VariantStorage string  `json:"variant_storage"`
	VariantRAM     string  `json:"variant_ram"`
	Quantity       int     `json:"quantity" binding:"required,min=1"`
	SellPrice      float64 `json:"sell_price" binding:"required,gt=0"`
	OrderSource    string  `json:"order_source"`
	AccountingType string  `json:"accounting_type"`
	ManagerID      *int    `json:"manager_id"`
	CreatedAt      *string `json:"created_at"`
}
