package controller

import (
	"tesla_server/model"
)

func CronJobs() {
	// checkDeviceCreditTicker := time.NewTicker(1 * time.Minute)
	// minuteTicker := time.NewTicker(10 * time.Minute)
	// // hourTicker := time.NewTicker(1 * time.Hour)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-checkDeviceCreditTicker.C:
	// 			go checkDeviceCredit()
	// 		case <-minuteTicker.C:
	// 			// go safeJob()
	// 			fmt.Println("run sec")
	// 		}
	// 	}
	// }()

	checkDeviceCredit()
}

func checkDeviceCredit() {
	var devices []model.Device
	devices, _ = model.GetDevicesByTeslaStream(0)
	var vins []string
	for _, device := range devices {
		tesla_stream := 0
		vins = nil
		vins = append(vins, device.Vin)
		result := ConnectDevice(vins, device.AccessToken)

		if result {
			tesla_stream = 2
		} else {
			tesla_stream = 1
		}
		model.UpdateDeviceTeslaStreambyVin(device.Vin, tesla_stream)
	}
}
