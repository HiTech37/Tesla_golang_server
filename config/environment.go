package config

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TeslaCredential struct {
	ClientID       string
	SecretKey      string
	CallbackUri    string
	DataScope      string
	RootDomain     string
	ServerDomain   string
	Certificate    string
	ProxyUri       string
	Port           int
	TlsCertificate string
	ClientCert     string
	ClientKey      string
}

func GetTeslaCredential() *TeslaCredential {
	var teslaCredential TeslaCredential

	var environment = "prod" // "local", "test", "prod"

	if environment == "local" {
		teslaCredential.CallbackUri = "http://localhost:3000/tesla_signup"
	} else if environment == "test" {
		teslaCredential.CallbackUri = "https://test.moovetrax.com/tesla_signup"

	} else if environment == "prod" {
		teslaCredential.CallbackUri = "https://moovetrax.com/tesla_signup"
	}

	var certificate = `-----BEGIN CERTIFICATE-----
MIICZDCCAcagAwIBAgITN/0lwafuTVzh2pQf4T5gnv301zAKBggqhkjOPQQDAjAh
MR8wHQYDVQQDDBZ0ZXNsYWFwaS5tb292ZXRyYXguY29tMB4XDTI1MDIxMDIyMTAw
M1oXDTM1MDIwODIyMTAwM1owITEfMB0GA1UEAwwWdGVzbGFhcGkubW9vdmV0cmF4
LmNvbTCBmzAQBgcqhkjOPQIBBgUrgQQAIwOBhgAEAVPy46GJ4HvaYaxvaiU96+I/
9t8lFvKpcMy2og4+WxQayIdm6cRaIIIQrTapgnuypbWlMgCIWXdrxR5dYXMUeKrc
ABneHjWmK+KT4v7Wr8nTlesTDRjSg4uMeE4pAIB77xf2R4e2nMOc2aEXgd4Jzq8/
DU5Pve+5lZfabeurvSSB+kOho4GZMIGWMB0GA1UdDgQWBBTOwiGhjpzgm/a8aO90
05bUsRuFlTAfBgNVHSMEGDAWgBTOwiGhjpzgm/a8aO9005bUsRuFlTAPBgNVHRMB
Af8EBTADAQH/MCEGA1UdEQQaMBiCFnRlc2xhYXBpLm1vb3ZldHJheC5jb20wEwYD
VR0lBAwwCgYIKwYBBQUHAwEwCwYDVR0PBAQDAgKMMAoGCCqGSM49BAMCA4GLADCB
hwJBQ8klYfiJSc2zb8hW9DgnkY1BeZQAUETqVSvxzqyOf2vIb3NXN8RSzJzDPVWt
j/O6Ztkt9qwA6CEjJRZWLOW1JsECQgHFjKUbZwTlFvJaFvYmjZbisi6RDqerI/V2
XaXV6SwBQ9IM0U1zgXCL4lf6cVUPaoTL8c/BzKAOyJ7RF/YPlXeoYg==
-----END CERTIFICATE-----` // cert.pem

	teslaCredential.ClientID = "60d97918-9b6b-4c92-88e3-ff9e9403239f"
	teslaCredential.SecretKey = "ta-secret.8p8Jz&Y^%n9FCCeE"
	teslaCredential.RootDomain = "https://moovetrax.com"
	teslaCredential.DataScope = "openid vehicle_device_data vehicle_location offline_access vehicle_cmds vehicle_charging_cmds"
	teslaCredential.Certificate = certificate
	teslaCredential.ServerDomain = "fleetapi.moovetrax.com"
	teslaCredential.Port = 8443
	teslaCredential.ProxyUri = "https://fleetapi.moovetrax.com:4443"

	return &teslaCredential
}

// Connect MySql DB

var dbInstance *gorm.DB
var once sync.Once

func InitDb() (*gorm.DB, error) {
	var err error
	// Initialize the DB connection only once
	once.Do(func() {
		dbInstance, err = connectDB()
		if err != nil {
			log.Fatalf("Failed to connect to database: %s", err)
		}
	})

	return dbInstance, err
}

func connectDB() (*gorm.DB, error) {

	DB_NAME := "gpsdb"
	DB_HOST := ""
	DB_USER := ""
	DB_PASSWORD := ""
	DB_PORT := "3306"

	var db_connect = "test" // "local", "test", "prod"

	switch db_connect {
	case "local":
		DB_HOST = "127.0.0.1"
		DB_USER = "root"
	case "test":
		DB_HOST = "172.31.42.68"
		DB_USER = "root"
		DB_PASSWORD = "342A$$$SD1232"
	case "prod":
		DB_HOST = "moovetrax-1.ckmhdbxyagjk.us-east-1.rds.amazonaws.com"
		DB_USER = "moovetrx"
		DB_PASSWORD = "342A$$$SD1232"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting raw database connection: %w", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging the database: %w", err)
	}

	return db, nil
}
