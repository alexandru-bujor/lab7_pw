package category

import (
    "MegaMobileBack/pkg/db"
)

func GetCategoryTemplateByCategoryID(categoryID int) (*CategoryTemplate, error) {
    var tmpl CategoryTemplate
    if err := db.DB.Where("category_id = ?", categoryID).First(&tmpl).Error; err != nil {
        return nil, err
    }
    return &tmpl, nil
}

func UpsertCategoryTemplate(tmpl *CategoryTemplate) error {
    // if exists, update; else create
    var existing CategoryTemplate
    if err := db.DB.Where("category_id = ?", tmpl.CategoryID).First(&existing).Error; err == nil {
        tmpl.ID = existing.ID
        return db.DB.Save(tmpl).Error
    }
    return db.DB.Create(tmpl).Error
}

func DeleteCategoryTemplateByCategoryID(categoryID int) error {
    return db.DB.Where("category_id = ?", categoryID).Delete(&CategoryTemplate{}).Error
}
