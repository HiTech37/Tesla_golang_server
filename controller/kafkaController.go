package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"tesla_server/model"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// TelemetryTesla represents the database entity
type TelemetryTesla struct {
	ID                uint      `gorm:"primaryKey"`
	LocationLatitude  float64   `json:"location_latitude"`
	LocationLongitude float64   `json:"location_longitude"`
	BatteryLevel      float64   `json:"battery_level"`
	Odometer          float64   `json:"odometer"`
	VehicleSpeed      float64   `json:"vehicle_speed"`
	Vin               string    `json:"vin"`
	CreatedAt         time.Time `json:"createdAt"`
}

// TelemetryData represents the incoming JSON structure
type TelemetryData struct {
	CreatedAt string `json:"createdAt"`
	Vin       string `json:"vin"`
	Data      []struct {
		Key   string `json:"key"`
		Value struct {
			DoubleValue   float64 `json:"doubleValue"`
			LocationValue struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"locationValue"`
		} `json:"value"`
	} `json:"data"`
}

func KafkaConsumer() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "54.161.235.213:9093",
		"group.id":          "telemetry",
		"auto.offset.reset": "earliest",
		"message.max.bytes": 524288000,
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	// Subscribe to the topic
	err = consumer.Subscribe("telemetry_V", nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	fmt.Println("Kafka consumer is running...")
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			var telemetryData TelemetryData
			err := json.Unmarshal(msg.Value, &telemetryData)
			if err != nil {
				log.Printf("Failed to parse JSON: %s", err)
				continue
			}

			// Parse datetime
			createdAt := time.Now()
			createdAt, err = time.Parse(time.RFC3339, telemetryData.CreatedAt)
			if err != nil {
				log.Printf("Failed to parse datetime: %s", err)
				continue
			}

			// Initialize variables to capture telemetry values
			var latitude, longitude, batteryLevel, odometer, vehicleSpeed float64
			var vin string

			// Extract data
			vin = telemetryData.Vin

			for _, item := range telemetryData.Data {
				switch item.Key {
				case "Location":
					latitude = item.Value.LocationValue.Latitude
					longitude = item.Value.LocationValue.Longitude
				case "BatteryLevel":
					batteryLevel = item.Value.DoubleValue
				case "Odometer":
					odometer = item.Value.DoubleValue
				case "VehicleSpeed":
					vehicleSpeed = item.Value.DoubleValue
				}
			}

			var device model.Device
			var position model.Position
			device.BatteryLevel = batteryLevel
			position.BatteryLevel = batteryLevel
			device.Latitude = latitude
			position.Latitude = latitude
			device.Longitude = longitude
			position.Longitude = longitude
			device.Odometer = odometer
			position.Odometer = odometer
			device.Speed = int(vehicleSpeed)
			position.Speed = int(vehicleSpeed)
			device.Vin = vin

			if latitude != 0 && longitude != 0 {
				fmt.Println("debug=>", device.Speed)
				fmt.Println("debug=>", position.Speed)
				err = model.UpdateDeviceByVin(device)
				if err != nil {
					fmt.Println(err)
				}

				position.DeviceTime = createdAt
				err = model.AddPositionInfo(position, vin)
				if err != nil {
					fmt.Println(err)
				}
			}

		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
