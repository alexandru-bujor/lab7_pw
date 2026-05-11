package lombard

import (
	"MegaMobileBack/pkg/db"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService() *Service {
	return &Service{DB: db.DB}
}

func (s *Service) Create(input CreateLombardRequestInput) (*LombardRequest, error) {
	request := &LombardRequest{
		CustomerName:  input.CustomerName,
		CustomerPhone: input.CustomerPhone,
		RequestType:   input.RequestType,
		Products:      input.Products,
		Status:        "new",
	}

	if request.RequestType == "" {
		request.RequestType = "electronics"
	}

	if err := s.DB.Create(request).Error; err != nil {
		fmt.Printf("❌ Error saving Lombard Request: %v\n", err)
		return nil, err
	}

	fmt.Printf("✅ Saved Lombard Request ID %d with %d products\n", request.ID, len(request.Products))

	return request, nil
}

func (s *Service) List() ([]LombardRequest, error) {
	var requests []LombardRequest
	if err := s.DB.Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (s *Service) GetByID(id int) (*LombardRequest, error) {
	var request LombardRequest
	if err := s.DB.First(&request, id).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

func (s *Service) UpdateStatus(id int, status string) (*LombardRequest, error) {
	var request LombardRequest
	if err := s.DB.First(&request, id).Error; err != nil {
		return nil, errors.New("lombard request not found")
	}

	validStatuses := map[string]bool{
		"new":      true,
		"old":      true,
		"seen":     true,
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

func (s *Service) Delete(id int) error {
	return s.DB.Delete(&LombardRequest{}, id).Error
}
