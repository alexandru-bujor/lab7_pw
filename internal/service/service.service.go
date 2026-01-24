package service

import (
	"MegaMobileBack/pkg/db"
)

func GetServicesByPart(part int) ([]Service, error) {
	var services []Service
	err := db.DB.Where("service_part = ?", part).Find(&services).Error
	return services, err
}

func GetAllServices() ([]Service, error) {
	var services []Service
	err := db.DB.Find(&services).Error
	return services, err
}

func GetServiceByID(id uint) (*Service, error) {
	var service Service
	err := db.DB.First(&service, id).Error
	return &service, err
}

func CreateService(s *Service) error {
	return db.DB.Create(s).Error
}

func UpdateService(id uint, s *Service) error {
	return db.DB.Model(&Service{}).Where("id = ?", id).Updates(s).Error
}

func DeleteService(id uint) error {
	return db.DB.Delete(&Service{}, id).Error
}
