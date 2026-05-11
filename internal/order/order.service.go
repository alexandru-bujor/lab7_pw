package order

import (
	"MegaMobileBack/internal/client"
	"MegaMobileBack/internal/inventory"
	"MegaMobileBack/internal/product"
	"MegaMobileBack/pkg/db"
	"errors"
	"fmt"
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

func (s *Service) List(status string, startDate string, endDate string) ([]Order, error) {
	var orders []Order
	// Select all fields from orders and use a join to get the manager name
	q := s.DB.Table("orders").
		Select("orders.*, CONCAT_WS(' ', users.first_name, users.last_name) as manager_name, products.brand_name as product_brand_name, products.model as product_model").
		Joins("LEFT JOIN users ON users.id = orders.manager_id").
		Joins("LEFT JOIN products ON products.id = orders.product_id")

	if status != "" {
		q = q.Where("orders.status = ?", status)
	}
	if startDate != "" {
		q = q.Where("DATE(orders.created_at) >= ?", startDate)
	}
	if endDate != "" {
		q = q.Where("DATE(orders.created_at) <= ?", endDate)
	}

	if err := q.Order("orders.created_at DESC").Find(&orders).Error; err != nil {
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

func (s *Service) UpdateStatus(id int, status string, managerID *int) (*Order, error) {
	var order Order
	if err := s.DB.First(&order, id).Error; err != nil {
		return nil, err
	}

	order.Status = status
	if managerID != nil {
		order.ManagerID = managerID
	}
	if err := s.DB.Save(&order).Error; err != nil {
		return nil, err
	}

	// Fetch the manager name for the response
	if order.ManagerID != nil {
		var managerName string
		s.DB.Table("users").Select("CONCAT_WS(' ', first_name, last_name)").Where("id = ?", *order.ManagerID).Scan(&managerName)
		order.ManagerName = managerName
	}

	return &order, nil
}

func (s *Service) CreateFromWebsite(input CreateOrderInput) (*Order, error) {
	if input.ClientEmail == "" {
		return nil, errors.New("client email is required")
	}
	if !strings.Contains(input.ClientEmail, "@") {
		return nil, errors.New("invalid email format")
	}

	clientSvc := client.NewService()
	clientSvc.DB = s.DB

	clientRecord, err := clientSvc.GetOrCreateClientByEmail(input.ClientEmail, input.ClientName, input.ClientPhone)
	if err != nil {
		return nil, err
	}

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
		ManagerID:      input.ManagerID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if input.CreatedAt != nil && *input.CreatedAt != "" {
		if t, err := time.Parse("2006-01-02", *input.CreatedAt); err == nil {
			// Maintain original time of day if it's already there? No, manual creates usually mean it was on that day.
			// Set it to 12:00 UTC or just keep as is (which will be 00:00:00).
			order.CreatedAt = t
		}
	}

	if order.OrderSource == "" {
		order.OrderSource = "website"
	}

	// If manual sale, make it visible in secondary panel by default
	if order.OrderSource == "manual" {
		order.IsVisibleInSecondary = true
	}

	if order.AccountingType == "" {
		order.AccountingType = "on_book"
	}

	if err := s.DB.Create(order).Error; err != nil {
		return nil, err
	}

	// Fetch the manager name if manager exists
	if order.ManagerID != nil {
		var managerName string
		s.DB.Table("users").Select("CONCAT_WS(' ', first_name, last_name)").Where("id = ?", *order.ManagerID).Scan(&managerName)
		order.ManagerName = managerName
	}

	return order, nil
}

func (s *Service) MarkAsSold(orderID int, inventoryItemID int, sellPrice *float64, managerID *int) (*Order, error) {
	var order Order
	if err := s.DB.First(&order, orderID).Error; err != nil {
		return nil, err
	}

	inventorySvc := inventory.NewService()
	inventorySvc.DB = s.DB

	var item inventory.InventoryItem
	if err := s.DB.First(&item, inventoryItemID).Error; err != nil {
		return nil, errors.New("inventory item not found")
	}

	if item.ProductID != order.ProductID {
		return nil, errors.New("inventory item does not match order product")
	}

	if item.Status == "sold" {
		return nil, errors.New("inventory item already sold")
	}

	var prod product.Product
	if err := s.DB.First(&prod, item.ProductID).Error; err != nil {
		return nil, err
	}

	now := time.Now()

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		item.Status = "sold"
		item.SoldAt = &now
		if order.ClientID != nil {
			item.ClientID = order.ClientID
		}
		item.OrderID = &order.ID
		if sellPrice != nil {
			item.SellPrice = *sellPrice
		}
		if err := tx.Save(&item).Error; err != nil {
			return fmt.Errorf("failed to save inventory item: %w", err)
		}

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
		if managerID != nil {
			order.ManagerID = managerID
		}

		if err := tx.Save(&order).Error; err != nil {
			return fmt.Errorf("failed to save order: %w", err)
		}

		if order.ClientID != nil {
			var clientRecord client.Client
			if err := tx.First(&clientRecord, *order.ClientID).Error; err == nil {
				clientRecord.OrderCount++
				clientRecord.TotalSpent += order.SellPrice
				clientRecord.LastOrderDate = &now
				if err := tx.Save(&clientRecord).Error; err != nil {
					return fmt.Errorf("failed to update client record: %w", err)
				}
			}
		}

		if prod.Stock > 0 {
			prod.Stock--
			if err := tx.Model(&product.Product{}).Where("id = ?", prod.ID).Update("stock", prod.Stock).Error; err != nil {
				return fmt.Errorf("failed to decrement product stock: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Fetch the manager name for the response
	if order.ManagerID != nil {
		var managerName string
		s.DB.Table("users").Select("CONCAT_WS(' ', first_name, last_name)").Where("id = ?", *order.ManagerID).Scan(&managerName)
		order.ManagerName = managerName
	}

	return &order, nil
}

func (s *Service) Delete(id int) error {
	return s.DB.Delete(&Order{}, id).Error
}

func (s *Service) GetSecondaryOrders() ([]Order, error) {
	var orders []Order
	q := s.DB.Table("orders").
		Select("orders.*, CONCAT_WS(' ', users.first_name, users.last_name) as manager_name, products.brand_name as product_brand_name, products.model as product_model").
		Joins("LEFT JOIN users ON users.id = orders.manager_id").
		Joins("LEFT JOIN products ON products.id = orders.product_id").
		Where("orders.is_visible_in_secondary = ?", true)

	if err := q.Order("orders.created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
