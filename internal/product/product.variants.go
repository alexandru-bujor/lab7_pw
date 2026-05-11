package product

import "MegaMobileBack/pkg/db"

func GetProductVariants(productID int) ([]Product, error) {
	var p Product
	if err := db.DB.First(&p, productID).Error; err != nil {
		return []Product{}, err
	}

	var variants []Product
	err := db.DB.
		Preload("Images").
		Where("brand_name = ? AND model = ? AND category_id = ? AND `condition` = ? AND is_published = ? AND (accounting_type = '' OR accounting_type = 'on_book')",
			p.BrandName, p.Model, p.CategoryID, p.Condition, true).
		Find(&variants).Error

	return variants, err
}
