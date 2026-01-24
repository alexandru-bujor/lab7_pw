package inventory

import (
	"MegaMobileBack/internal/product"
	"MegaMobileBack/pkg/db"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService() *Service { return &Service{DB: db.DB} }

func (s *Service) List(section string, status string) ([]InventoryItem, error) {
	var items []InventoryItem
	q := s.DB.Model(&InventoryItem{})
	if section != "" { q = q.Where("section = ?", section) }
	if status != "" { q = q.Where("status = ?", status) }
	if err := q.Order("updated_at DESC").Find(&items).Error; err != nil { return nil, err }
	return items, nil
}

func (s *Service) Create(item *InventoryItem) error {
	return s.DB.Create(item).Error
}

func (s *Service) Update(id int, payload map[string]any) (*InventoryItem, error) {
	var item InventoryItem
	if err := s.DB.First(&item, id).Error; err != nil { return nil, err }
	if err := s.DB.Model(&item).Updates(payload).Error; err != nil { return nil, err }
	return &item, nil
}

func (s *Service) Delete(id int) error { return s.DB.Delete(&InventoryItem{}, id).Error }

func (s *Service) AssignClient(id int, clientID int) (*InventoryItem, error) {
	var item InventoryItem
	if err := s.DB.First(&item, id).Error; err != nil { return nil, err }
	now := time.Now()
	item.ClientID = &clientID
	item.AssignedAt = &now
	item.Status = "reserved"
	if err := s.DB.Save(&item).Error; err != nil { return nil, err }
	return &item, nil
}

func (s *Service) MarkSold(id int, clientID *int, sellPrice *float64) (*InventoryItem, error) {
	var item InventoryItem
	if err := s.DB.First(&item, id).Error; err != nil { return nil, err }
	
	// Get product details
	var prod product.Product
	if err := s.DB.Preload("Variants").First(&prod, item.ProductID).Error; err != nil { return nil, err }
	
	now := time.Now()
	item.Status = "sold"
	item.SoldAt = &now
	if clientID != nil { item.ClientID = clientID }
	if sellPrice != nil { item.SellPrice = *sellPrice }
	
	// Start transaction
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Save inventory item as sold
		if err := tx.Save(&item).Error; err != nil { return err }
		
		// Decrement stock - check if variant exists
		hasVariant := item.VariantColor != "" || item.VariantStorage != "" || item.VariantRAM != ""
		if hasVariant {
			// Find and decrement variant stock
			var variant product.ProductVariant
			err := tx.Where("product_id = ? AND color = ? AND storage = ? AND ram = ?", 
				item.ProductID, item.VariantColor, item.VariantStorage, item.VariantRAM).
				First(&variant).Error
			if err == nil && variant.Stock > 0 {
				variant.Stock--
				if err := tx.Save(&variant).Error; err != nil { return err }
			}
		} else {
			// Decrement product stock
			if prod.Stock > 0 {
				prod.Stock--
				if err := tx.Model(&product.Product{}).Where("id = ?", prod.ID).Update("stock", prod.Stock).Error; err != nil { 
					return err 
				}
			}
		}
		
		// Create order record directly (no package import needed)
		itemID := item.ID
		orderRecord := map[string]interface{}{
			"client_id":         item.ClientID,
			"client_name":       "",
			"client_email":      "",
			"client_phone":      "",
			"inventory_item_id": &itemID,
			"product_id":        item.ProductID,
			"product_title":     prod.Title,
			"variant_color":     item.VariantColor,
			"variant_storage":   item.VariantStorage,
			"variant_ram":       item.VariantRAM,
			"quantity":          1,
			"sell_price":        item.SellPrice,
			"buy_price":         item.BuyPrice,
			"profit":            item.SellPrice - item.BuyPrice,
			"status":            "completed",
			"order_source":      "manual",
			"sold_at":           &now,
		}
		if err := tx.Table("orders").Create(orderRecord).Error; err != nil { return err }
		
		return nil
	})
	
	if err != nil { return nil, err }
	return &item, nil
}
// Sections CRUD
func (s *Service) ListSections() ([]InventorySection, error) {
    var sections []InventorySection
    if err := s.DB.Order("name ASC").Find(&sections).Error; err != nil { return nil, err }
    return sections, nil
}

func (s *Service) CreateSection(name string) (*InventorySection, error) {
    section := &InventorySection{Name: name}
    if err := s.DB.Create(section).Error; err != nil { return nil, err }
    return section, nil
}

func (s *Service) UpdateSection(id int, name string) (*InventorySection, error) {
    var sec InventorySection
    if err := s.DB.First(&sec, id).Error; err != nil { return nil, err }
    sec.Name = name
    if err := s.DB.Save(&sec).Error; err != nil { return nil, err }
    return &sec, nil
}

func (s *Service) DeleteSection(id int) error {
    return s.DB.Delete(&InventorySection{}, id).Error
}