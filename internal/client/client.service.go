package client

import (
	"MegaMobileBack/pkg/db"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Service handles client business logic
type Service struct {
	DB *gorm.DB
}

// NewService creates a new client service
func NewService() *Service {
	return &Service{DB: db.DB}
}

// GetAllClients retrieves all clients with their order history
func (s *Service) GetAllClients() ([]Client, error) {
	var clients []Client
	err := s.DB.Preload("Orders").Find(&clients).Error
	return clients, err
}

// GetClientByID retrieves a client by ID with order history
func (s *Service) GetClientByID(id int) (*Client, error) {
	var client Client
	err := s.DB.Preload("Orders.Items").First(&client, id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// GetClientByEmail retrieves a client by email
func (s *Service) GetClientByEmail(email string) (*Client, error) {
	var client Client
	// Use Find instead of First to avoid logging "record not found" as error
	err := s.DB.Where("email = ?", email).Limit(1).Find(&client).Error
	if err != nil {
		return nil, err
	}
	// If ID is 0, no record was found
	if client.ID == 0 {
		return nil, nil
	}
	return &client, nil
}

// CreateClient creates a new client
func (s *Service) CreateClient(input CreateClientInput) (*Client, error) {
	// Check if client already exists
	existing, err := s.GetClientByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("client with this email already exists")
	}

	client := Client{
		Name:       input.Name,
		Email:      input.Email,
		Phone:      input.Phone,
		OrderCount: 0,
		TotalSpent: 0,
	}

	err = s.DB.Create(&client).Error
	if err != nil {
		return nil, err
	}

	return &client, nil
}

// UpdateClient updates an existing client
func (s *Service) UpdateClient(id int, input UpdateClientInput) (*Client, error) {
	client, err := s.GetClientByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		client.Name = input.Name
	}
	if input.Email != "" {
		client.Email = input.Email
	}
	if input.Phone != "" {
		client.Phone = input.Phone
	}

	err = s.DB.Save(client).Error
	if err != nil {
		return nil, err
	}

	return client, nil
}

// DeleteClient deletes a client
func (s *Service) DeleteClient(id int) error {
	return s.DB.Delete(&Client{}, id).Error
}

// RecordOrder creates a new order for a client and updates their stats
func (s *Service) RecordOrder(clientID int, orderInput struct {
	OrderNumber string  `json:"order_number" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required"`
	Status      string  `json:"status"`
	Notes       string  `json:"notes"`
}) (*ClientOrder, error) {
	// Start transaction
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get client
	var client Client
	if err := tx.First(&client, clientID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create order
	now := time.Now()
	order := ClientOrder{
		ClientID:    clientID,
		OrderNumber: orderInput.OrderNumber,
		TotalAmount: orderInput.TotalAmount,
		Status:      orderInput.Status,
		OrderDate:   now,
		Notes:       orderInput.Notes,
	}

	if order.Status == "" {
		order.Status = "pending"
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update client stats
	client.OrderCount++
	client.TotalSpent += orderInput.TotalAmount
	client.LastOrderDate = &now

	if err := tx.Save(&client).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrCreateClientByEmail gets existing client or creates new one (for auto-registration from orders)
func (s *Service) GetOrCreateClientByEmail(email, name, phone string) (*Client, error) {
	// Try to get existing client
	client, err := s.GetClientByEmail(email)
	if err != nil {
		return nil, err
	}

	// If client exists, return it
	if client != nil {
		return client, nil
	}

	// Otherwise create new client
	return s.CreateClient(CreateClientInput{
		Name:  name,
		Email: email,
		Phone: phone,
	})
}

// GetSecondaryClients returns only clients visible in secondary panel
func (s *Service) GetSecondaryClients() ([]Client, error) {
	var clients []Client
	err := s.DB.Where("is_visible_in_secondary = ?", true).Preload("Orders").Find(&clients).Error
	return clients, err
}