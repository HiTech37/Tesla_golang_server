package model

import (
	"tesla_server/config"
)

type Device struct {
	ID            uint   `json:"id"`
	VehicleID     string `json:"vehicleId"`
	Vin           string `json:"vin"`
	AccessToken   string `json:"accessToken"`
	RefreshToken  string `json:"refreshToken"`
	IsPaid        int    `json:"is_paid"`
	BillingSource string `json:"billing_source"`
	MonthlyCost   string `json:"monthly_cost"`
	Credit        string `json:"credit"`
	TeslaStream   int    `json:"tesla_stream"`
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
