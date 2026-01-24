
package servicecategory

import (
	"MegaMobileBack/pkg/db"
)

func UpdateServiceCategory(id string, sc *ServiceCategory) error {
	return db.DB.Model(&ServiceCategory{}).Where("id = ?", id).Updates(sc).Error
}

func DeleteServiceCategory(id string) error {
	return db.DB.Delete(&ServiceCategory{}, id).Error
}

func CreateServiceCategory(sc *ServiceCategory) error {
	return db.DB.Create(sc).Error
}

func GetAllServiceCategories() ([]ServiceCategory, error) {
	var categories []ServiceCategory
	err := db.DB.Find(&categories).Error
	return categories, err
}
