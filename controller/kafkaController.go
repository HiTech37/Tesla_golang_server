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
	ChargeState       string    `json:"charge_state"`
	CreatedAt         time.Time `json:"createdAt"`
}

type TelemetryData struct {
	CreatedAt string `json:"createdAt"`
	Data      []struct {
		Key   string `json:"key"`
		Value struct {
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
		"bootstrap.servers": "fleetapi.moovetrax.com:9093",
		"group.id":          "telemetry",
		// "client.id":         "telemetry-service",
		"auto.offset.reset": "earliest",
		"message.max.bytes": 524288000,
		// 	"security.protocol":        "SSL",
		// 	"ssl.ca.location":          "/home/ec2-user/key_store/fleetapi.moovetrax.com/cert.pem",
		// 	"ssl.certificate.location": "/home/ec2-user/key_store/fleetapi.moovetrax.com/fullchain_fixed.pem",
		// 	"ssl.key.location":         "/home/ec2-user/key_store/fleetapi.moovetrax.com/privkey.pem",
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

	// Continuously consume messages
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

			var latitude, longitude float64
			var chargeState string

			// Extract Location and ChargeState data
			for _, item := range telemetryData.Data {
				if item.Key == "Location" {
					latitude = item.Value.LocationValue.Latitude
					longitude = item.Value.LocationValue.Longitude
				} else if item.Key == "ChargeState" {
					chargeStateBytes, err := json.Marshal(item.Value)
					if err == nil {
						chargeState = string(chargeStateBytes)
					}
				}
			}

			// Save data to the database if latitude and longitude are present
			if latitude != 0 && longitude != 0 {
				telemetry := TelemetryTesla{
					LocationLatitude:  latitude,
					LocationLongitude: longitude,
					ChargeState:       chargeState,
					CreatedAt:         createdAt,
				}

				fmt.Println("=>", telemetry)

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
