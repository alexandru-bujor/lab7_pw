package product

import (
	"MegaMobileBack/pkg/db"
)

// GetProductVariants returns all products in the same variant group
func GetProductVariants(productID int) ([]Product, error) {
	var product Product
	if err := db.DB.First(&product, productID).Error; err != nil {
		return []Product{}, err
	}

	// If the requested product is off-book, don't return variants (clients shouldn't see off-book products)
	if product.AccountingType == "off_book" {
		return []Product{}, nil
	}

	// If this product has no variant group, return just itself (if it's on-book)
	if product.VariantGroupID == nil {
		return []Product{product}, nil
	}

	// Get all products in the same variant group, but only on-book ones
	var variants []Product
	err := db.DB.
		Preload("Images").
		Where("(variant_group_id = ? OR id = ?) AND (accounting_type IS NULL OR accounting_type = '' OR accounting_type = 'on_book')", 
			*product.VariantGroupID, *product.VariantGroupID).
		Find(&variants).Error

	return variants, err
}
