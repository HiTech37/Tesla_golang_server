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
	TestServerUri  string
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

	var certificate = `-----BEGIN CERTIFICATE-----
MIICJzCCAYigAwIBAgIUOYyl49rR2nfXfGUZnA8K+KSUtEYwCgYIKoZIzj0EAwIw
FDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI1MDMxMzIwMjU1NVoXDTM1MDMxMTIw
MjU1NVowFDESMBAGA1UEAwwJbG9jYWxob3N0MIGbMBAGByqGSM49AgEGBSuBBAAj
A4GGAAQBNP/F4j9zAsvsF9e2rQRwWRUVaaXHlXWBtmgAwfJgf/fgXNhyR/Hqw1tl
2jws7aY6Nrj79d44mmeCbnG/4Vr/n9UBvCHusZ0LPcsRBr5VSmj53xGLOA2RBfFM
vWuipjdUqJVteDdrRpGSH5aSF2m+Iteiv3aMWVY1/K2bPoWWJIEGbG+jdTBzMB0G
A1UdDgQWBBQyxhJlAf97/eqxEIG1+HAD8hGs1zAfBgNVHSMEGDAWgBQyxhJlAf97
/eqxEIG1+HAD8hGs1zAPBgNVHRMBAf8EBTADAQH/MBMGA1UdJQQMMAoGCCsGAQUF
BwMBMAsGA1UdDwQEAwICjDAKBggqhkjOPQQDAgOBjAAwgYgCQgGcXOAPybKGNAyk
YRQiSCTC1h2dupa+UmujRpouyUz0v+qyXfv5EPznpQ8eFQUDd9NY856MTTupIQUv
CnWgVJ88aAJCAKZjRSSrgw4lAHXr6qyCyGcCoVr8SZPuXuqqDk7i+BjJ0XEHePPg
M4FqRaFSnjY1EadGi2UPCEqk8fARdUwA0EmN
-----END CERTIFICATE-----` // cert.pem
	teslaCredential.CallbackUri = "https://fleetapi.moovetrax.com/device_signup"
	teslaCredential.ClientID = "60d97918-9b6b-4c92-88e3-ff9e9403239f"
	teslaCredential.SecretKey = "ta-secret.8p8Jz&Y^%n9FCCeE"
	teslaCredential.RootDomain = "https://moovetrax.com"
	teslaCredential.DataScope = "openid%20vehicle_device_data%20vehicle_location%20offline_access%20vehicle_cmds%20vehicle_charging_cmds"
	teslaCredential.Certificate = certificate
	teslaCredential.ServerDomain = "fleetapi.moovetrax.com"
	teslaCredential.Port = 8443
	teslaCredential.ProxyUri = "https://fleetapi.moovetrax.com:4443"
	// teslaCredential.ProxyUri = "https://localhost:4443"

	teslaCredential.TestServerUri = "https://test.moovetrax.com:8088"

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
