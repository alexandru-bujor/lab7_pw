package product

import (
	"MegaMobileBack/pkg/db"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GetAllProducts() ([]Product, error) {
	var products []Product
	err := db.DB.
		Preload("Images").
		Find(&products).Error

	return products, err
}

func AddProduct(p *Product) error {

	if len(p.Images) > 0 {
		seen := map[string]bool{}
		dedup := make([]ProductImage, 0, len(p.Images))
		for _, img := range p.Images {
			if img.URL == "" {
				continue
			}
			if seen[img.URL] {
				continue
			}
			seen[img.URL] = true
			dedup = append(dedup, img)
		}
		p.Images = dedup
	}

	// Set a temporary unique slug to avoid unique constraint violation before we have an ID
	p.Slug = fmt.Sprintf("temp-%d", time.Now().UnixNano())

	err := db.DB.Create(p).Error
	if err != nil {
		return err
	}

	p.Slug = GenerateSlug(p.Title, p.ID)
	return db.DB.Model(p).Update("slug", p.Slug).Error
}

func DeleteProduct(id int) error {
	tx := db.DB.Begin()

	if err := tx.Where("product_id = ?", id).Delete(&ProductImage{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&Product{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func UpdateProduct(id int, p *Product) error {
	p.ID = id

	var existingProduct Product
	if err := db.DB.Preload("Images").First(&existingProduct, id).Error; err == nil {
		if p.Title == "" {
			p.Title = existingProduct.Title
		}
		if p.Slug == "" {
			p.Slug = existingProduct.Slug
		}
		if p.Description == "" {
			p.Description = existingProduct.Description
		}
		if p.CategoryID == 0 {
			p.CategoryID = existingProduct.CategoryID
		}
		if p.BrandName == "" {
			p.BrandName = existingProduct.BrandName
		}
		if p.BrandID == nil {
			p.BrandID = existingProduct.BrandID
		}
		if p.DeviceType == "" {
			p.DeviceType = existingProduct.DeviceType
		}
		if p.Specs == nil {
			p.Specs = existingProduct.Specs
		}
		if p.Images == nil {
			p.Images = existingProduct.Images
		}
		if p.AccountingType == "" {
			p.AccountingType = existingProduct.AccountingType
		}
		if existingProduct.Title != p.Title {
			p.Slug = GenerateSlug(p.Title, p.ID)
		}
		p.CreatedAt = existingProduct.CreatedAt
	}
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{FullSaveAssociations: false}).Save(p).Error; err != nil {
			return err
		}

		if err := tx.Where("product_id = ?", p.ID).Delete(&ProductImage{}).Error; err != nil {
			return err
		}

		seen := map[string]bool{}
		var toInsert []ProductImage
		for _, img := range p.Images {
			if img.URL == "" {
				continue
			}
			if seen[img.URL] {
				continue
			}
			seen[img.URL] = true
			toInsert = append(toInsert, ProductImage{
				ProductID: p.ID,
				URL:       img.URL,
				AltText:   img.AltText,
				IsMain:    img.IsMain,
			})
		}

		if len(toInsert) > 0 {
			if err := tx.Create(&toInsert).Error; err != nil {
				return err
			}
		}

		return tx.Preload("Images").First(p, p.ID).Error
	})
}

func GetProductBySlug(slug string) (*Product, error) {
	var p Product
	err := db.DB.
		Preload("Images").
		Where("slug = ? AND (accounting_type = '' OR accounting_type = 'on_book')", slug).
		First(&p).Error

	return &p, err
}
