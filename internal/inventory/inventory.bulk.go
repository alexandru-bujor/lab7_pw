package inventory

import (
	"MegaMobileBack/internal/product"
	"MegaMobileBack/pkg/db"
)

func SyncAllProductsToInventory(defaultSection string) (int, error) {
	if defaultSection == "" {
		defaultSection = "Apple SH"
	}

	var products []product.Product
	if err := db.DB.Find(&products).Error; err != nil {
		return 0, err
	}

	createdCount := 0
	for _, p := range products {
		var existingCount int64
		db.DB.Model(&InventoryItem{}).Where("product_id = ?", p.ID).Count(&existingCount)

		missing := p.Stock - int(existingCount)

		for i := 0; i < missing; i++ {
			item := InventoryItem{
				ProductID:      p.ID,
				Section:        defaultSection,
				VariantColor:   p.VariantColor,
				VariantStorage: p.VariantStorage,
				VariantRAM:     p.VariantRAM,
				BuyPrice:       p.Price,
				SellPrice:      p.Price,
				Status:         "in_stock",
			}
			if err := db.DB.Create(&item).Error; err != nil {
				continue
			}
			createdCount++
		}
	}

	return createdCount, nil
}
