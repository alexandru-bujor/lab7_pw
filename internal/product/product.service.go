package product

import (
	"MegaMobileBack/pkg/db"

	"gorm.io/gorm"
)

func GetAllProducts() ([]Product, error) {
	var products []Product
	err := db.DB.
		Preload("Images").
		Preload("Variants").
		Find(&products).Error

	return products, err
}

func AddProduct(p *Product) error {
	// Check if there are existing products with same title and category
	var existingProducts []Product
	db.DB.Where("title = ? AND category_id = ?", p.Title, p.CategoryID).Find(&existingProducts)
	
	// If existing products found, use their variant_group_id
	if len(existingProducts) > 0 {
		if existingProducts[0].VariantGroupID != nil {
			p.VariantGroupID = existingProducts[0].VariantGroupID
		} else {
			// First product didn't have variant_group_id, use its ID
			groupID := existingProducts[0].ID
			p.VariantGroupID = &groupID
			// Update the first product to have variant_group_id
			db.DB.Model(&existingProducts[0]).Update("variant_group_id", groupID)
		}
	}
	
	// Deduplicate images by URL before create
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

	// Create the product
	err := db.DB.Create(p).Error
	if err != nil {
		return err
	}
	
	// If this is the first product and no variant group was set, use its own ID
	if p.VariantGroupID == nil {
		p.VariantGroupID = &p.ID
		db.DB.Model(p).Update("variant_group_id", p.ID)
	}
	
	// Update with generated slug
	p.Slug = GenerateSlug(p.Title, p.ID)
	return db.DB.Model(p).Update("slug", p.Slug).Error
}

func DeleteProduct(id int) error {
	return db.DB.Delete(&Product{}, id).Error
}

func UpdateProduct(id int, p *Product) error {
	// Ensure we update the correct row and persist nested relations
	p.ID = id
	
	// Regenerate slug if title changed
	var existingProduct Product
	if err := db.DB.First(&existingProduct, id).Error; err == nil {
		if existingProduct.Title != p.Title {
			p.Slug = GenerateSlug(p.Title, p.ID)
		}
	}
	// Replace images set to avoid duplicates: remove existing and insert new
	// Run in a transaction to keep data consistent
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Save basic fields on product first
		if err := tx.Session(&gorm.Session{FullSaveAssociations: false}).Save(p).Error; err != nil {
			return err
		}

		// Clear existing images for this product
		if err := tx.Where("product_id = ?", p.ID).Delete(&ProductImage{}).Error; err != nil {
			return err
		}

		// Deduplicate images by URL before re-inserting
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

		// Replace variants set similarly (optional): if provided, clear and insert
		if p.Variants != nil {
			if err := tx.Where("product_id = ?", p.ID).Delete(&ProductVariant{}).Error; err != nil {
				return err
			}
			if len(p.Variants) > 0 {
				var variantsToInsert []ProductVariant
				for _, v := range p.Variants {
					variantsToInsert = append(variantsToInsert, ProductVariant{
						ProductID: p.ID,
						Storage:   v.Storage,
						RAM:       v.RAM,
						Color:     v.Color,
						Price:     v.Price,
					})
				}
				if err := tx.Create(&variantsToInsert).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
