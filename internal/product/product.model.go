package product

import (
	"time"
	"gorm.io/datatypes"
)

type Product struct {
	ID              int       `gorm:"column=id;primaryKey;autoIncrement" json:"id"`
	Title           string    `gorm:"column=title" json:"title"`
	Slug            string    `gorm:"column=slug;unique" json:"slug"`
	Description     string    `gorm:"column=description" json:"description"`
	CategoryID      int       `gorm:"column=category_id" json:"category_id"`
	BrandID         *int      `gorm:"column=brand_id" json:"brand_id,omitempty"`
	BrandName       string    `gorm:"column=brand_name" json:"brand_name"`
	Price           float64   `gorm:"column=price" json:"price"`
	Stock           int       `gorm:"column=stock" json:"stock"`
	DeviceType      string    `gorm:"column=device_type" json:"device_type"`
	Condition       string                 `gorm:"column=condition" json:"condition"`
	Specs           datatypes.JSONMap      `gorm:"column=specs;type:json" json:"specs"`
	VariantGroupID  *int      `gorm:"column=variant_group_id" json:"variant_group_id,omitempty"`
	VariantColor    string    `gorm:"column=variant_color" json:"variant_color"`
	VariantStorage  string    `gorm:"column=variant_storage" json:"variant_storage"`
	VariantRAM      string    `gorm:"column=variant_ram" json:"variant_ram"`
	AccountingType  string    `gorm:"column=accounting_type;default:on_book" json:"accounting_type"` // on_book | off_book
	IsVisibleInSecondary bool  `gorm:"column=is_visible_in_secondary;default:false" json:"is_visible_in_secondary"`
	CreatedAt       time.Time `gorm:"column=created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column=updated_at" json:"updated_at"`

	Images   []ProductImage   `gorm:"foreignKey:ProductID" json:"images"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID" json:"variants"`
}

func (Product) TableName() string { return "products" }

type ProductImage struct {
	ID        int    `gorm:"column=id;primaryKey;autoIncrement" json:"id"`
	ProductID int    `gorm:"column=product_id" json:"product_id"`
	URL       string `gorm:"column=url" json:"url"`
	AltText   string `gorm:"column=alt_text" json:"alt"`
	IsMain    bool   `gorm:"column=is_main" json:"is_main"`
}

func (ProductImage) TableName() string { return "product_images" }

type ProductVariant struct {
	ID        int     `gorm:"column=id;primaryKey;autoIncrement" json:"id"`
	ProductID int     `gorm:"column=product_id" json:"product_id"`
	Storage   string  `gorm:"column=storage" json:"storage"`
	RAM       string  `gorm:"column=ram" json:"ram"`
	Color     string  `gorm:"column=color" json:"color"`
	Price     float64 `gorm:"column=price" json:"price"`
	Stock     int     `gorm:"column=stock" json:"stock"`
}

func (ProductVariant) TableName() string { return "product_variants" }
