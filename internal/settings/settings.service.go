package settings

import (
	"MegaMobileBack/internal/utils"
	"MegaMobileBack/pkg/db"
	"errors"
)

func GetSetting(key string) (string, error) {
	var setting Setting
	result := db.DB.Where("key = ?", key).First(&setting)
	if result.Error != nil {
		return "", errors.New("setting not found")
	}
	return setting.Value, nil
}

func SetSetting(key, value string) error {
	var setting Setting
	result := db.DB.Where("key = ?", key).First(&setting)

	if result.Error != nil {
		setting = Setting{
			Key:   key,
			Value: value,
		}
		return db.DB.Create(&setting).Error
	}

	setting.Value = value
	return db.DB.Save(&setting).Error
}

func GetSecondaryPanelPassword() (string, error) {
	password, err := GetSetting(SettingKeySecondaryPanelPassword)
	if err != nil {
		defaultPassword := "admin123"
		hashedPassword := utils.HashPasswordSHA256(defaultPassword)
		SetSetting(SettingKeySecondaryPanelPassword, hashedPassword)
		return defaultPassword, nil
	}
	return password, nil
}

func SetSecondaryPanelPassword(plainPassword string) error {
	hashedPassword := utils.HashPasswordSHA256(plainPassword)
	return SetSetting(SettingKeySecondaryPanelPassword, hashedPassword)
}

func VerifySecondaryPanelPassword(plainPassword string) (bool, error) {
	hashedPassword, err := GetSetting(SettingKeySecondaryPanelPassword)
	if err != nil {
		defaultPassword := "admin123"
		return utils.CheckPasswordSHA256(plainPassword, utils.HashPasswordSHA256(defaultPassword)), nil
	}
	return utils.CheckPasswordSHA256(plainPassword, hashedPassword), nil
}
