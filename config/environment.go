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
		teslaCredential.CallbackUri = "https://fleetapi.moovetrax.com/device_signup"
	}

	var certificate = `-----BEGIN CERTIFICATE-----
MIICZjCCAcegAwIBAgIUQf7+xIcYqtENFZ6Hb4aXyrvS94gwCgYIKoZIzj0EAwIw
ITEfMB0GA1UEAwwWZmxlZXRhcGkubW9vdmV0cmF4LmNvbTAeFw0yNTAyMjQxMjQz
MTdaFw0zNTAyMjIxMjQzMTdaMCExHzAdBgNVBAMMFmZsZWV0YXBpLm1vb3ZldHJh
eC5jb20wgZswEAYHKoZIzj0CAQYFK4EEACMDgYYABAFCuqni5ZkWiySsgitfTESY
AbMQSfHVvaoI3Tp/0y1Z3VO81jg5WVuLDoIP7+4lrhtgNRYy2ujsMn6NyWGxXbpx
MwEFymz6SFRk0Z5KcSJ//drfH3T7ySx8z8AwTWuX3b1X2VM8kPGuWLCRf023L1sY
qyp3f6UxkYkvOZU4u+vUjVoQM6OBmTCBljAdBgNVHQ4EFgQUnsfMHx3uw54TzUeX
0aR5tb0VN6gwHwYDVR0jBBgwFoAUnsfMHx3uw54TzUeX0aR5tb0VN6gwDwYDVR0T
AQH/BAUwAwEB/zAhBgNVHREEGjAYghZmbGVldGFwaS5tb292ZXRyYXguY29tMBMG
A1UdJQQMMAoGCCsGAQUFBwMBMAsGA1UdDwQEAwICjDAKBggqhkjOPQQDAgOBjAAw
gYgCQgCpkKGBIzUJ7rLW2iDYPF54XPnPCjGiNEkENaaDxXhw7fNrpyLcvBjJUc2+
/IeCYeHEpPrfhsd/7CJJL1r/SGWWUQJCAP6wePV7Pkg6VT/Tlq9rPypEBwshOD9R
mDx+/ZKSmVS87QfXse+m5yh/1/O1EdpQXcvyN+tAjhW7tPWDaHbAuDN9
-----END CERTIFICATE-----` // cert.pem

	teslaCredential.ClientID = "60d97918-9b6b-4c92-88e3-ff9e9403239f"
	teslaCredential.SecretKey = "ta-secret.8p8Jz&Y^%n9FCCeE"
	teslaCredential.RootDomain = "https://moovetrax.com"
	teslaCredential.DataScope = "openid%20vehicle_device_data%20vehicle_location%20offline_access%20vehicle_cmds%20vehicle_charging_cmds"
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
