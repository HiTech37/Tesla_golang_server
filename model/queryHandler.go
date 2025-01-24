package model

import (
	"tesla_server/config"

	"gorm.io/gorm"
)

var db *gorm.DB

func UpdateDeviceAuthTokensbyVin(accessToken string, refreshToken string, vin string) error {
	db, err := config.InitDb()
	if err != nil {
		return err
	}

	result := db.Exec("UPDATE devices SET accessToken = ?, refreshToken = ? WHERE vin = ?", accessToken, refreshToken, vin)
	if result.Error != nil {
		return err
	}

	return nil
}
