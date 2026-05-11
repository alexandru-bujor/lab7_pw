package client

import (
	"time"
)

type Client struct {
	ID                   int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name                 string     `gorm:"column:name;not null" json:"name"`
	Email                string     `gorm:"column:email;unique;not null" json:"email"`
	Phone                string     `gorm:"column:phone;not null" json:"phone"`
	OrderCount           int        `gorm:"column:order_count;default:0" json:"order_count"`
	TotalSpent           float64    `gorm:"column:total_spent;default:0" json:"total_spent"`
	LastOrderDate        *time.Time `gorm:"column:last_order_date" json:"last_order_date,omitempty"`
	IsVisibleInSecondary bool       `gorm:"column=is_visible_in_secondary;default:false" json:"is_visible_in_secondary"`
	CreatedAt            time.Time  `gorm:"column=created_at;autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"column=updated_at;autoUpdateTime" json:"updated_at"`

	Orders []ClientOrder `gorm:"foreignKey:ClientID" json:"orders,omitempty"`
}

func (Client) TableName() string { return "clients" }

type ClientOrder struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ClientID    int       `gorm:"column:client_id;not null" json:"client_id"`
	OrderNumber string    `gorm:"column:order_number;unique;not null" json:"order_number"`
	TotalAmount float64   `gorm:"column:total_amount;not null" json:"total_amount"`
	Status      string    `gorm:"column:status;default:'pending'" json:"status"`
	OrderDate   time.Time `gorm:"column:order_date;autoCreateTime" json:"order_date"`
	Notes       string    `gorm:"column:notes" json:"notes,omitempty"`

	Client Client            `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	Items  []ClientOrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func (ClientOrder) TableName() string { return "client_orders" }

type ClientOrderItem struct {
	ID          int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderID     int     `gorm:"column:order_id;not null" json:"order_id"`
	ProductID   int     `gorm:"column:product_id" json:"product_id"`
	ProductName string  `gorm:"column:product_name;not null" json:"product_name"`
	Quantity    int     `gorm:"column:quantity;not null" json:"quantity"`
	Price       float64 `gorm:"column:price;not null" json:"price"`
	Subtotal    float64 `gorm:"column:subtotal;not null" json:"subtotal"`
}

func (ClientOrderItem) TableName() string { return "client_order_items" }

type CreateClientInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
}

type UpdateClientInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
