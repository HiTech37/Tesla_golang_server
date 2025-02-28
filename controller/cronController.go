package controller

import (
	"fmt"
	"tesla_server/model"
	"time"
)

func CronJobs() {
	checkDeviceCreditTicker := time.NewTicker(1 * time.Minute)
	hourTicker := time.NewTicker(1 * time.Minute)
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
	var devices []model.Device
	devices, _ = model.GetDevicesByTeslaStream(0)
	var vins []string
	for _, device := range devices {
		vins = nil
		vins = append(vins, device.Vin)
		tesla_stream := ConnectDevice(vins, device.AccessToken, device.RefreshToken)
		model.UpdateDeviceTeslaStreambyVin(device.Vin, tesla_stream)
		if tesla_stream == 1 {
			UpdateUnSupportedDeviceInfo(device.Vin, device.AccessToken)
		}
	}
}

func handleUnsupportedDevice() {
	var devices []model.Device
	devices, _ = model.GetDevicesByTeslaStream(1)
	for _, device := range devices {
		fmt.Println("debug3=>", device.Vin)
		err := UpdateUnSupportedDeviceInfo(device.Vin, device.AccessToken)
		if err != nil {
			fmt.Println("debug2=>", err)
		}
	}
}
