package model

import (
	"fmt"
	"io"
	"net"
	"net/http"
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
	PrevOdometer  float64   `json:"prevOdometer" gorm:"column:prev_od"`
	Odometer      float64   `json:"odometer"`
	Status        string    `json:"status"`
	LastPosition  time.Time `json:"lastPosition" gorm:"column:lastPosition"`
	LastConnect   time.Time `json:"lastConnect" gorm:"column:lastConnect"`
}

type Position struct {
	DeviceId     int       `json:"deviceId" gorm:"column:deviceId"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Speed        int       `json:"speed"`
	BatteryLevel float64   `json:"batteryLevel" gorm:"column:mt2v_dc_volt"`
	Odometer     float64   `json:"odometer"`
	DeviceTime   time.Time `json:"deviceTime" gorm:"column:deviceTime"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updatedAt"`
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
	result := db.Where(`
    deviceType = ? 
    AND (
        (credit <= 2 * monthly_cost) 
        AND 
        (billing_source = ? AND credit <= 0) 
        OR 
        (billing_source != ? AND is_paid = ?)
    ) 
    AND tesla_stream = ?`,
		"tesla",
		"escrow",
		"escrow",
		1,
		tesla_stream,
	).Find(&devices)

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

	// Set the PrevOdometer to the current Odometer value
	deviceInfo.PrevOdometer = existingDevice.Odometer

	// Update the record
	if err := db.Model(&existingDevice).Updates(deviceInfo).Error; err != nil {
		return err
	}

	db.Model(&existingDevice).Select("Speed").Updates(Device{Speed: deviceInfo.Speed})

	return nil

}

func UpdateDeviceByVin(deviceInfo Device, isSpeedUpdated bool) error {
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

	if isSpeedUpdated {
		db.Model(&existingDevice).Select("Speed").Updates(Device{Speed: deviceInfo.Speed})
	}

	return nil

}

func AddPositionInfo(vin string, createdAt time.Time) error {
	db, err := config.InitDb()
	if err != nil {
		return err
	}

	var device Device
	if err := db.Where("vin = ?", vin).First(&device).Error; err != nil {
		return err
	}

	var position Position
	if device.Latitude != 0 && device.Longitude != 0 {
		position.Latitude = device.Latitude
		position.Longitude = device.Longitude
		position.DeviceId = int(device.ID)
		position.BatteryLevel = device.BatteryLevel
		position.Odometer = device.Odometer
		position.Speed = device.Speed
		position.DeviceTime = createdAt
		position.CreatedAt = time.Now()
		position.UpdatedAt = time.Now()
		if err := db.Create(&position).Error; err != nil {
			return err
		}
	}

	return nil
}

func UpdateHandshake(vin string) error {
	db, err := config.InitDb()
	if err != nil {
		return err
	}

	upload_ip := getPrivateIP()
	upload_public_ip := getPublicIP()
	device_ip := ""
	device_port := ""

	query := `INSERT INTO handshakes (gps_id, upload_ip, upload_public_ip, device_ip, device_port, handshake_cnt, createdAt, updatedAt)
              VALUES (?, ?, ?, ?, ?, 0, NOW(), NOW())
              ON DUPLICATE KEY UPDATE handshake_cnt = handshake_cnt + 1, updatedAt = NOW();`

	result := db.Exec(query, vin, upload_ip, upload_public_ip, device_ip, device_port)
	if result.Error != nil {
		return fmt.Errorf("failed to execute query: %v", result.Error)
	}

	return nil

}

func getPrivateIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func getPublicIP() string {
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(ip)
}
