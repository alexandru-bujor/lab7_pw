package client

import (
	"MegaMobileBack/pkg/db"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService() *Service {
	return &Service{DB: db.DB}
}

func (s *Service) GetAllClients() ([]Client, error) {
	var clients []Client
	err := s.DB.Preload("Orders").Find(&clients).Error
	return clients, err
}

func (s *Service) GetClientByID(id int) (*Client, error) {
	var client Client
	err := s.DB.Preload("Orders.Items").First(&client, id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (s *Service) GetClientByEmail(email string) (*Client, error) {
	var client Client
	err := s.DB.Where("email = ?", email).Limit(1).Find(&client).Error
	if err != nil {
		return nil, err
	}
	if client.ID == 0 {
		return nil, nil
	}
	return &client, nil
}

func (s *Service) CreateClient(input CreateClientInput) (*Client, error) {
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

func (s *Service) DeleteClient(id int) error {
	return s.DB.Delete(&Client{}, id).Error
}

func (s *Service) RecordOrder(clientID int, orderInput struct {
	OrderNumber string  `json:"order_number" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required"`
	Status      string  `json:"status"`
	Notes       string  `json:"notes"`
}) (*ClientOrder, error) {
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var client Client
	if err := tx.First(&client, clientID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

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

func (s *Service) GetOrCreateClientByEmail(email, name, phone string) (*Client, error) {
	client, err := s.GetClientByEmail(email)
	if err != nil {
		return nil, err
	}

	if client != nil {
		return client, nil
	}

	return s.CreateClient(CreateClientInput{
		Name:  name,
		Email: email,
		Phone: phone,
	})
}

func (s *Service) GetSecondaryClients() ([]Client, error) {
	var clients []Client
	err := s.DB.Where("is_visible_in_secondary = ?", true).Preload("Orders").Find(&clients).Error
	return clients, err
}
