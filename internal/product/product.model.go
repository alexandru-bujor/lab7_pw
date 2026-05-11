package product

import (
	"time"

	"gorm.io/datatypes"
)

type Product struct {
	ID                   int               `gorm:"column=id;primaryKey;autoIncrement" json:"id"`
	Title                string            `gorm:"column=title" json:"title"`
	Slug                 string            `gorm:"column=slug;unique" json:"slug"`
	Description          string            `gorm:"column=description" json:"description"`
	CategoryID           int               `gorm:"column=category_id" json:"category_id"`
	BrandID              *int              `gorm:"column=brand_id" json:"brand_id,omitempty"`
	BrandName            string            `gorm:"column=brand_name" json:"brand_name"`
	Price                float64           `gorm:"column=price" json:"price"`
	OldPrice             float64           `gorm:"column=old_price;default:0" json:"old_price"`
	Stock                int               `gorm:"column=stock" json:"stock"`
	DeviceType           string            `gorm:"column=device_type" json:"device_type"`
	Condition            string            `gorm:"column=condition" json:"condition"`
	Model                string            `gorm:"column=model" json:"model"`
	ProductCode          string            `gorm:"column=product_code" json:"product_code"`
	BatteryHealth        string            `gorm:"column=battery_health" json:"battery_health"`
	IsRecommended        bool              `gorm:"column=is_recommended;default:false" json:"is_recommended"`
	IsNewArrival         bool              `gorm:"column=is_new_arrival;default:false" json:"is_new_arrival"`
	Specs                datatypes.JSONMap `gorm:"column=specs;type:json" json:"specs"`
	VariantColor         string            `gorm:"column=variant_color" json:"variant_color"`
	VariantStorage       string            `gorm:"column=variant_storage" json:"variant_storage"`
	VariantRAM           string            `gorm:"column=variant_ram" json:"variant_ram"`
	AccountingType       string            `gorm:"column=accounting_type;default:on_book" json:"accounting_type"`
	IsVisibleInSecondary bool              `gorm:"column=is_visible_in_secondary;default:false" json:"is_visible_in_secondary"`
	IsPublished          bool              `gorm:"column=is_published;default:true" json:"is_published"`
	CreatedAt            time.Time         `gorm:"column=created_at" json:"created_at"`
	UpdatedAt            time.Time         `gorm:"column=updated_at" json:"updated_at"`

	Images []ProductImage `gorm:"foreignKey:ProductID" json:"images"`
}

func (Product) TableName() string { return "products" }

type ProductImage struct {
	ID        int    `gorm:"column=id;primaryKey;autoIncrement" json:"id"`
	ProductID int    `gorm:"column=product_id" json:"product_id"`
	URL       string `gorm:"column=url" json:"url"`
	AltText   string `gorm:"column=alt_text" json:"alt_text"`
	IsMain    bool   `gorm:"column=is_main" json:"is_main"`
}

func (ProductImage) TableName() string { return "product_images" }
