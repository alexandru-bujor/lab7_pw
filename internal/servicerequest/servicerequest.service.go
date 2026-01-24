package servicerequest

import (
	"MegaMobileBack/pkg/db"
	"errors"

	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService() *Service {
	return &Service{DB: db.DB}
}

// Create creates a new service request
func (s *Service) Create(input CreateServiceRequestInput) (*ServiceRequest, error) {
	request := &ServiceRequest{
		CustomerName:      input.CustomerName,
		CustomerPhone:     input.CustomerPhone,
		DeviceMake:        input.DeviceMake,
		DeviceModel:       input.DeviceModel,
		ProblemDescription: input.ProblemDescription,
		Status:            "new",
	}

	if err := s.DB.Create(request).Error; err != nil {
		return nil, err
	}

	return request, nil
}

// List returns all service requests
func (s *Service) List() ([]ServiceRequest, error) {
	var requests []ServiceRequest
	if err := s.DB.Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// GetByID returns a specific service request
func (s *Service) GetByID(id int) (*ServiceRequest, error) {
	var request ServiceRequest
	if err := s.DB.First(&request, id).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

// UpdateStatus updates the status of a service request
func (s *Service) UpdateStatus(id int, status string) (*ServiceRequest, error) {
	var request ServiceRequest
	if err := s.DB.First(&request, id).Error; err != nil {
		return nil, errors.New("service request not found")
	}

	// Validate status
	validStatuses := map[string]bool{
		"new":  true,
		"old":  true,
		"seen": true,
		// Keep old statuses for backward compatibility
		"pending":  true,
		"approved": true,
		"accepted": true,
		"rejected": true,
		"declined": true,
	}
	if !validStatuses[status] {
		return nil, errors.New("invalid status")
	}

	request.Status = status
	if err := s.DB.Save(&request).Error; err != nil {
		return nil, err
	}

	return &request, nil
}

// Delete deletes a service request
func (s *Service) Delete(id int) error {
	return s.DB.Delete(&ServiceRequest{}, id).Error
}

