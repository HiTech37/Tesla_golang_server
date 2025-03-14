package controller

import (
	"fmt"
	"tesla_server/model"
	"time"
)

func CronJobs() {
	checkDeviceCreditTicker := time.NewTicker(1 * time.Minute)
	hourTicker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-checkDeviceCreditTicker.C:
				go checkDeviceCredit()
			case <-hourTicker.C:
				handleUnsupportedDevice()
			}
		}
	}()
}

func checkDeviceCredit() {
	fmt.Println("Checking the deivce credit......")
	var devices []model.Device
	devices, _ = model.GetDevicesByTeslaStream(0)
	var vins []string
	for _, device := range devices {
		vins = nil
		vins = append(vins, device.Vin)
		tesla_stream := ConnectDevice(vins, device.AccessToken, device.RefreshToken)
		err := model.UpdateDeviceTeslaStreambyVin(device.Vin, tesla_stream)
		if err != nil {
			fmt.Println(err)
		}
		if tesla_stream >= 1 {
			err := UpdateUnSupportedDeviceInfo(device.Vin, device.AccessToken, device.RefreshToken)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func handleUnsupportedDevice() {
	fmt.Println("Handling Unsupported device......")
	var devices []model.Device
	devices, _ = model.GetDevicesByTeslaStream(1)
	for _, device := range devices {
		fmt.Println("devices=>", device.Vin)
		err := UpdateUnSupportedDeviceInfo(device.Vin, device.AccessToken, device.RefreshToken)
		if err != nil {
			fmt.Println(err)
		}
	}
}
