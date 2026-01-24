package settings

import (
	"MegaMobileBack/internal/utils"
	"MegaMobileBack/pkg/db"
	"errors"
)

// GetSetting retrieves a setting by key
func GetSetting(key string) (string, error) {
	var setting Setting
	result := db.DB.Where("key = ?", key).First(&setting)
	if result.Error != nil {
		return "", errors.New("setting not found")
	}
	return setting.Value, nil
}

// SetSetting creates or updates a setting
func SetSetting(key, value string) error {
	var setting Setting
	result := db.DB.Where("key = ?", key).First(&setting)
	
	if result.Error != nil {
		// Setting doesn't exist, create it
		setting = Setting{
			Key:   key,
			Value: value,
		}
		return db.DB.Create(&setting).Error
	}
	
	// Update existing setting
	setting.Value = value
	return db.DB.Save(&setting).Error
}

// GetSecondaryPanelPassword retrieves the secondary panel password
func GetSecondaryPanelPassword() (string, error) {
	password, err := GetSetting(SettingKeySecondaryPanelPassword)
	if err != nil {
		// Default password if not set
		defaultPassword := "admin123"
		hashedPassword := utils.HashPasswordSHA256(defaultPassword)
		SetSetting(SettingKeySecondaryPanelPassword, hashedPassword)
		return defaultPassword, nil
	}
	return password, nil
}

// SetSecondaryPanelPassword sets the secondary panel password (expects plain text, will hash)
func SetSecondaryPanelPassword(plainPassword string) error {
	hashedPassword := utils.HashPasswordSHA256(plainPassword)
	return SetSetting(SettingKeySecondaryPanelPassword, hashedPassword)
}

// VerifySecondaryPanelPassword verifies the provided password against stored hash
func VerifySecondaryPanelPassword(plainPassword string) (bool, error) {
	hashedPassword, err := GetSetting(SettingKeySecondaryPanelPassword)
	if err != nil {
		// If setting doesn't exist, check against default
		defaultPassword := "admin123"
		return utils.CheckPasswordSHA256(plainPassword, utils.HashPasswordSHA256(defaultPassword)), nil
	}
	return utils.CheckPasswordSHA256(plainPassword, hashedPassword), nil
}
