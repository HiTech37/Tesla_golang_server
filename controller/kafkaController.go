package controller

import (
	"encoding/json"
	"fmt"
	"log"
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
	// Initialize Kafka consumer
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

			// Extract data
			vin = telemetryData.Vin // Capture VIN

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

			// Save data to the database if latitude and longitude are present
			if latitude != 0 && longitude != 0 {
				telemetry := TelemetryTesla{
					LocationLatitude:  latitude,
					LocationLongitude: longitude,
					BatteryLevel:      batteryLevel,
					Odometer:          odometer,
					VehicleSpeed:      vehicleSpeed,
					Vin:               vin,
					CreatedAt:         createdAt,
				}

				test(telemetry)

				// Uncomment and adjust the following lines if you're using GORM for database operations
				// if err := db.WithContext(context.Background()).Create(&telemetry).Error; err != nil {
				// 	log.Printf("Failed to save to database: %s", err)
				// } else {
				// 	fmt.Printf("Saved to database: %+v\n", telemetry)
				// }
			}
		} else {
			// Log consumer errors
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func test(telemetry TelemetryTesla) {

}
