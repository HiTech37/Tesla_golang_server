package model

import (
	"tesla_server/config"

	"gorm.io/gorm"
)

type TeslaDeviceStatus struct {
	DeviceID    int     `gorm:"not null" json:"deviceid"`
	Vin         string  `gorm:"not null" json:"vin"`
	Fireware    string  `gorm:"not null" json:"firmware"`
	IsSupported bool    `gorm:"not null" json:"isSupported"`
	IsConnected bool    `gorm:"not null" json:"isConnected"`
	Data        *string `gorm:"type:text" json:"data"`
	gorm.Model
}

type TeslaDeviceStatusRepo struct {
	Db *gorm.DB
}

func CreateTeslaDeviceStatusTable() (*TeslaDeviceStatusRepo, error) {
	db, err := config.InitDb()
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&TeslaDeviceStatus{})
	if err != nil {
		return nil, err
	}

	return &TeslaDeviceStatusRepo{Db: db}, nil
}

func CreateTeslaDeviceStatus(teslaDeviceStatus *TeslaDeviceStatus) error {
	teslaDeviceStatusRepo, err := CreateTeslaDeviceStatusTable()
	if err != nil {
		return err
	}

	err = teslaDeviceStatusRepo.Db.Create(teslaDeviceStatus).Error
	if err != nil {
		return err
	}
	return nil
}
