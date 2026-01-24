package make

import (
	"MegaMobileBack/pkg/db"
)

func CreateMake(m *Make) error {
	return db.DB.Create(m).Error
}

func GetAllMakes() ([]Make, error) {
	var makes []Make
	err := db.DB.Preload("Models").Find(&makes).Error
	return makes, err
}

func UpdateMake(id uint, m *Make) error {
	return db.DB.Model(&Make{}).Where("id = ?", id).Updates(m).Error
}

func DeleteMake(id uint) error {
	return db.DB.Delete(&Make{}, id).Error
}

func CreateModel(model *Model) error {
	return db.DB.Create(model).Error
}

func UpdateModel(id uint, model *Model) error {
	return db.DB.Model(&Model{}).Where("id = ?", id).Updates(model).Error
}

func DeleteModel(id uint) error {
	return db.DB.Delete(&Model{}, id).Error
}
