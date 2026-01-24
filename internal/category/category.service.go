package category

import "MegaMobileBack/pkg/db"

func GetAllCategories() ([]Category, error) {
	var categories []Category
	err := db.DB.Find(&categories).Error
	return categories, err
}

func AddCategory(c *Category) error {
	return db.DB.Create(c).Error
}

func UpdateCategory(c *Category) error {
	return db.DB.Save(c).Error
}

func DeleteCategory(id int) error {
	return db.DB.Delete(&Category{}, id).Error
}
