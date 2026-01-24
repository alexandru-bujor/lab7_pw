package lombard

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

// Create creates a new lombard request
func (s *Service) Create(input CreateLombardRequestInput) (*LombardRequest, error) {
	request := &LombardRequest{
		CustomerName:  input.CustomerName,
		CustomerPhone: input.CustomerPhone,
		RequestType:   input.RequestType,
		Products:      input.Products,
		Status:        "new",
	}

	if request.RequestType == "" {
		request.RequestType = "electronics" // default
	}

	if err := s.DB.Create(request).Error; err != nil {
		return nil, err
	}

	return request, nil
}

// List returns all lombard requests
func (s *Service) List() ([]LombardRequest, error) {
	var requests []LombardRequest
	if err := s.DB.Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// GetByID returns a specific lombard request
func (s *Service) GetByID(id int) (*LombardRequest, error) {
	var request LombardRequest
	if err := s.DB.First(&request, id).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

// UpdateStatus updates the status of a lombard request
func (s *Service) UpdateStatus(id int, status string) (*LombardRequest, error) {
	var request LombardRequest
	if err := s.DB.First(&request, id).Error; err != nil {
		return nil, errors.New("lombard request not found")
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

// Delete deletes a lombard request
func (s *Service) Delete(id int) error {
	return s.DB.Delete(&LombardRequest{}, id).Error
}

