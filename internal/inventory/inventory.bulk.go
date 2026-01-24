package inventory

import (
	"MegaMobileBack/internal/product"
	"MegaMobileBack/pkg/db"
	"log"
)

// SyncAllProductsToInventory creates inventory items for all products and their variants
// that don't already exist in inventory. Items are created with status "in_stock" and
// can be filled with serial, code, source later.
func SyncAllProductsToInventory(defaultSection string) (int, error) {
	if defaultSection == "" {
		defaultSection = "Apple SH" // fallback
	}

	var products []product.Product
	if err := db.DB.Preload("Variants").Find(&products).Error; err != nil {
		return 0, err
	}

	createdCount := 0
	for _, p := range products {
		// If product has variants, create items per variant
		if len(p.Variants) > 0 {
			for _, v := range p.Variants {
				// Check how many items exist for this product+variant combo
				var existingCount int64
				db.DB.Model(&InventoryItem{}).Where("product_id = ? AND variant_color = ? AND variant_storage = ? AND variant_ram = ?",
					p.ID, v.Color, v.Storage, v.RAM).Count(&existingCount)
				
				stock := v.Stock
				missing := stock - int(existingCount)
				
				// Create missing items up to variant stock
				for i := 0; i < missing; i++ {
					item := InventoryItem{
						ProductID:      p.ID,
						Section:        defaultSection,
						VariantColor:   v.Color,
						VariantStorage: v.Storage,
						VariantRAM:     v.RAM,
						BuyPrice:       v.Price,
						SellPrice:      v.Price,
						Status:         "in_stock",
					}
					if err := db.DB.Create(&item).Error; err != nil {
						log.Printf("Failed to create inventory item for product %d variant: %v", p.ID, err)
						continue
					}
					createdCount++
				}
			}
		} else {
			// No variants: use product-level stock
			var existingCount int64
			db.DB.Model(&InventoryItem{}).Where("product_id = ?", p.ID).Count(&existingCount)
			
			stock := p.Stock
			missing := stock - int(existingCount)
			
			// Create missing items up to product stock
			for i := 0; i < missing; i++ {
				item := InventoryItem{
					ProductID: p.ID,
					Section:   defaultSection,
					BuyPrice:  p.Price,
					SellPrice: p.Price,
					Status:    "in_stock",
				}
				if err := db.DB.Create(&item).Error; err != nil {
					log.Printf("Failed to create inventory item for product %d: %v", p.ID, err)
					continue
				}
				createdCount++
			}
		}
	}

	return createdCount, nil
}
