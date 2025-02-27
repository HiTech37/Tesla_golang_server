package controller

import (
	"fmt"
	"tesla_server/model"
	"time"
)

func CronJobs() {
	checkDeviceCreditTicker := time.NewTicker(1 * time.Minute)
	minuteTicker := time.NewTicker(10 * time.Minute)
	// hourTicker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-checkDeviceCreditTicker.C:
				go checkDeviceCredit()
			case <-minuteTicker.C:
				// go safeJob()
				fmt.Println("run sec")
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
		fmt.Println("deviceData=>", device)
		result := ConnectDevice(vins, device.AccessToken, device.RefreshToken)

		model.UpdateDeviceTeslaStreambyVin(device.Vin, result)
	}
}
