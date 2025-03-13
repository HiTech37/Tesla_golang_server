package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"tesla_server/model"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// TelemetryData represents the incoming JSON structure
type TelemetryData struct {
	CreatedAt string `json:"createdAt"`
	Vin       string `json:"vin"`
	Data      []struct {
		Key   string `json:"key"`
		Value struct {
			DoubleValue   float64 `json:"doubleValue,omitempty"`
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
			createdAt, err := time.Parse(time.RFC3339, telemetryData.CreatedAt)
			if err != nil {
				log.Printf("Failed to parse datetime: %s", err)
				continue
			}

			// Initialize variables to capture telemetry values
			var latitude, longitude, batteryLevel, odometer, vehicleSpeed float64
			var vin string

			var isSpeedUpdated bool = false
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
					isSpeedUpdated = true
					vehicleSpeed = item.Value.DoubleValue
				}
			}

			var device model.Device
			device.BatteryLevel = batteryLevel
			device.Latitude = latitude
			device.Longitude = longitude
			device.Odometer = odometer
			device.Speed = int(vehicleSpeed)
			device.Vin = vin

			err = model.UpdateDeviceByVin(device, isSpeedUpdated)
			if err != nil {
				fmt.Println(err)
			}

			err = model.AddPositionInfo(vin, createdAt)
			if err != nil {
				fmt.Println(err)
			}

			err = model.UpdateHandshake(vin)
			if err != nil {
				fmt.Println(err)
			}

		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
