package order

import (
	"MegaMobileBack/internal/client"
	"MegaMobileBack/internal/inventory"
	"MegaMobileBack/internal/product"
	"MegaMobileBack/pkg/db"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService() *Service { return &Service{DB: db.DB} }

func (s *Service) Create(order *Order) error {
	return s.DB.Create(order).Error
}

func (s *Service) List(status string) ([]Order, error) {
	var orders []Order
	q := s.DB.Model(&Order{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *Service) GetByID(id int) (*Order, error) {
	var order Order
	if err := s.DB.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *Service) GetByClientID(clientID int) ([]Order, error) {
	var orders []Order
	if err := s.DB.Where("client_id = ?", clientID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *Service) UpdateStatus(id int, status string) (*Order, error) {
	var order Order
	if err := s.DB.First(&order, id).Error; err != nil {
		return nil, err
	}
	
	order.Status = status
	if err := s.DB.Save(&order).Error; err != nil {
		return nil, err
	}
	
	return &order, nil
}

// CreateFromWebsite creates an order from website/client submission
func (s *Service) CreateFromWebsite(input CreateOrderInput) (*Order, error) {
	// Basic email format validation (less strict than Go's email validator)
	if input.ClientEmail == "" {
		return nil, errors.New("client email is required")
	}
	if !strings.Contains(input.ClientEmail, "@") {
		return nil, errors.New("invalid email format")
	}
	
	// Get or create client
	clientSvc := client.NewService()
	clientSvc.DB = s.DB
	
	clientRecord, err := clientSvc.GetOrCreateClientByEmail(input.ClientEmail, input.ClientName, input.ClientPhone)
	if err != nil {
		return nil, err
	}
	
	// Create order
	order := &Order{
		ClientID:       &clientRecord.ID,
		ClientName:     input.ClientName,
		ClientEmail:    input.ClientEmail,
		ClientPhone:    input.ClientPhone,
		ProductID:      input.ProductID,
		ProductTitle:   input.ProductTitle,
		VariantColor:   input.VariantColor,
		VariantStorage: input.VariantStorage,
		VariantRAM:     input.VariantRAM,
		Quantity:       input.Quantity,
		SellPrice:      input.SellPrice,
		Status:         "pending",
		OrderSource:    input.OrderSource,
		AccountingType: input.AccountingType,
	}
	
	if order.OrderSource == "" {
		order.OrderSource = "website"
	}
	
	if order.AccountingType == "" {
		order.AccountingType = "on_book"
	}
	
	if err := s.DB.Create(order).Error; err != nil {
		return nil, err
	}
	
	return order, nil
}

// MarkAsSold marks an order as sold and links it to an inventory item
func (s *Service) MarkAsSold(orderID int, inventoryItemID int, sellPrice *float64) (*Order, error) {
	var order Order
	if err := s.DB.First(&order, orderID).Error; err != nil {
		return nil, err
	}
	
	// Get inventory item
	inventorySvc := inventory.NewService()
	inventorySvc.DB = s.DB
	
	var item inventory.InventoryItem
	if err := s.DB.First(&item, inventoryItemID).Error; err != nil {
		return nil, errors.New("inventory item not found")
	}
	
	// Verify inventory item matches order
	if item.ProductID != order.ProductID {
		return nil, errors.New("inventory item does not match order product")
	}
	
	// Check if already sold
	if item.Status == "sold" {
		return nil, errors.New("inventory item already sold")
	}
	
	// Get product details for buy price
	var prod product.Product
	if err := s.DB.Preload("Variants").First(&prod, item.ProductID).Error; err != nil {
		return nil, err
	}
	
	now := time.Now()
	
	// Start transaction
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Update inventory item as sold
		item.Status = "sold"
		item.SoldAt = &now
		if order.ClientID != nil {
			item.ClientID = order.ClientID
		}
		if sellPrice != nil {
			item.SellPrice = *sellPrice
		}
		if err := tx.Save(&item).Error; err != nil {
			return err
		}
		
		// Update order
		order.InventoryItemID = &item.ID
		order.BuyPrice = item.BuyPrice
		if sellPrice != nil {
			order.SellPrice = *sellPrice
		} else {
			order.SellPrice = item.SellPrice
		}
		order.Profit = order.SellPrice - order.BuyPrice
		order.Status = "completed"
		order.SoldAt = &now
		
		if err := tx.Save(&order).Error; err != nil {
			return err
		}
		
		// Update client stats
		if order.ClientID != nil {
			var clientRecord client.Client
			if err := tx.First(&clientRecord, *order.ClientID).Error; err == nil {
				clientRecord.OrderCount++
				clientRecord.TotalSpent += order.SellPrice
				clientRecord.LastOrderDate = &now
				if err := tx.Save(&clientRecord).Error; err != nil {
					return err
				}
			}
		}
		
		// Decrement stock
		hasVariant := item.VariantColor != "" || item.VariantStorage != "" || item.VariantRAM != ""
		if hasVariant {
			var variant product.ProductVariant
			err := tx.Where("product_id = ? AND color = ? AND storage = ? AND ram = ?",
				item.ProductID, item.VariantColor, item.VariantStorage, item.VariantRAM).
				First(&variant).Error
			if err == nil && variant.Stock > 0 {
				variant.Stock--
				if err := tx.Save(&variant).Error; err != nil {
					return err
				}
			}
		} else {
			if prod.Stock > 0 {
				prod.Stock--
				if err := tx.Model(&product.Product{}).Where("id = ?", prod.ID).Update("stock", prod.Stock).Error; err != nil {
					return err
				}
			}
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &order, nil
}

func (s *Service) Delete(id int) error {
	return s.DB.Delete(&Order{}, id).Error
}

// GetSecondaryOrders returns only orders visible in secondary panel
func (s *Service) GetSecondaryOrders() ([]Order, error) {
	var orders []Order
	if err := s.DB.Where("is_visible_in_secondary = ?", true).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}