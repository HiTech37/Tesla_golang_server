package model

import (
	"fmt"
	"tesla_server/config"
	"time"
)

type Device struct {
	ID            uint      `json:"id"`
	VehicleID     string    `json:"vehicleId"`
	Vin           string    `json:"vin"`
	AccessToken   string    `json:"accessToken" gorm:"column:accessToken"`
	RefreshToken  string    `json:"refreshToken" gorm:"column:refreshToken"`
	IsPaid        int       `json:"is_paid"`
	BillingSource string    `json:"billing_source"`
	MonthlyCost   float64   `json:"monthly_cost"`
	Credit        float64   `json:"credit"`
	TeslaStream   int       `json:"tesla_stream"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Speed         int       `json:"speed"`
	BatteryLevel  float64   `json:"batteryLevel" gorm:"column:mt2v_dc_volt"`
	PrevOdometer  float64   `json:"prevOdometer" gorm:"prev_od"`
	Odometer      float64   `json:"odometer"`
	Status        string    `json:"status"`
	LastPosition  time.Time `json:"lastPosition" gorm:"column:lastPosition"`
	LastConnect   time.Time `json:"lastConnect" gorm:"column:lastConnect"`
}

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

func UpdateDeviceTeslaStreambyVin(vin string, tesla_stream int) error {
	db, err := config.InitDb()
	if err != nil {
		return err
	}

	result := db.Exec("UPDATE devices SET tesla_stream = ? WHERE vin = ?", tesla_stream, vin)
	if result.Error != nil {
		return err
	}

	return nil
}

func GetDevicesByTeslaStream(tesla_stream int) ([]Device, error) {
	db, err := config.InitDb()
	if err != nil {
		return nil, err
	}

	var devices []Device
	result := db.Where(" deviceType = ? AND tesla_stream = ? AND ((credit > 2 * monthly_cost) OR (billing_source = ? AND credit <= 0) OR (is_paid = ?))",
		"tesla",
		tesla_stream,
		"escrow",
		0).Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}

	return devices, nil
}

func GetDeviceByVin(vin string) ([]Device, error) {
	var devices []Device
	db, err := config.InitDb()
	if err != nil {
		return devices, err
	}
	result := db.Where("vin = ?", vin)
	if result.Error != nil {
		return nil, result.Error
	}

	return devices, nil
}

func UpdateDeviceInfoByVin(deviceInfo Device) error {
	fmt.Println(deviceInfo)
	db, err := config.InitDb()
	if err != nil {
		return err
	}

	// Get the existing device record by VIN
	var existingDevice Device
	if err := db.Where("vin = ?", deviceInfo.Vin).First(&existingDevice).Error; err != nil {
		return err
	}

	// Update only the fields that are set (non-zero values)
	deviceInfo.LastConnect = time.Now()
	deviceInfo.LastPosition = time.Now()
	deviceInfo.Status = "online"

	// Set the PrevOdometer to the current Odometer value
	deviceInfo.PrevOdometer = existingDevice.Odometer

	// Update the record
	if err := db.Model(&existingDevice).Updates(deviceInfo).Error; err != nil {
		return err
	}

	return nil

}
