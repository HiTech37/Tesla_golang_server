package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gorm.io/gorm"
)

// TelemetryTesla struct (same as the model)
type TelemetryTesla struct {
	ID                uint      `gorm:"primaryKey"`
	LocationLatitude  float64   `json:"location_latitude"`
	LocationLongitude float64   `json:"location_longitude"`
	ChargeState       string    `json:"charge_state"`
	CreatedAt         time.Time `json:"createdAt"`
}

var db *gorm.DB

// func initDB() {
// 	var err error
// 	dsn := "host=localhost user=postgres password=your_password dbname=telemetry_db port=5432 sslmode=disable TimeZone=Asia/Kolkata"
// 	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Failed to connect to the database: %v", err)
// 	}

// 	// Auto Migrate the schema
// 	db.AutoMigrate(&TelemetryTesla{})
// }

// func saveData(data TelemetryTesla) error {
// 	result := db.Create(&data)
// 	return result.Error
// }

func KafkaConsumer() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9093",
		"group.id":          "telemetry",
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	topic := "telemetry_V"
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	fmt.Println("Listening for messages on topic:", topic)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			var telemetryData map[string]interface{}
			err := json.Unmarshal(msg.Value, &telemetryData)
			if err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			// Extracting data from JSON
			datestring, ok := telemetryData["createdAt"].(string)
			if !ok {
				log.Println("createdAt not found or invalid")
				continue
			}
			createdAt, err := time.Parse(time.RFC3339, datestring)
			if err != nil {
				log.Printf("Error parsing createdAt: %v", err)
				continue
			}

			dataList, ok := telemetryData["data"].([]interface{})
			if !ok {
				log.Println("data list not found or invalid")
				continue
			}

			var locationValue map[string]interface{}
			var chargeState string

			// Extract Location and ChargeState
			for _, item := range dataList {
				fmt.Println("Data1=>", item)
				dataMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				if dataMap["key"] == "Location" {
					if val, exists := dataMap["value"].(map[string]interface{}); exists {
						locationValue = val["locationValue"].(map[string]interface{})
					}
				}

				if dataMap["key"] == "ChargeState" {
					if val, exists := dataMap["value"].(string); exists {
						chargeState = val
					}
				}
			}

			if locationValue != nil {
				latitude, latOk := locationValue["latitude"].(float64)
				longitude, longOk := locationValue["longitude"].(float64)

				if latOk && longOk {
					// Prepare data to be saved
					data := TelemetryTesla{
						LocationLatitude:  latitude,
						LocationLongitude: longitude,
						ChargeState:       chargeState,
						CreatedAt:         createdAt,
					}

					fmt.Println("Data2=>", data)
					// Save data
					// err := saveData(data)
					if err != nil {
						log.Printf("Error saving data: %v", err)
					} else {
						log.Println("Data saved successfully")
					}
				}
			}
		} else {
			// Error handling
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
