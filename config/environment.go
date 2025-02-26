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
MIICZjCCAcegAwIBAgIUb3SeJMsfe5V7ER9m3ZXtp0nltDgwCgYIKoZIzj0EAwIw
ITEfMB0GA1UEAwwWZmxlZXRhcGkubW9vdmV0cmF4LmNvbTAeFw0yNTAyMjUxNDE4
MzRaFw0zNTAyMjMxNDE4MzRaMCExHzAdBgNVBAMMFmZsZWV0YXBpLm1vb3ZldHJh
eC5jb20wgZswEAYHKoZIzj0CAQYFK4EEACMDgYYABAGwVQ+SijE6W6rdBRyvvaTK
WM3YKIrTFe0FDo7JlPvJbE5xUtQaBHL72HwkiAClqiy5ca2h7NKs5phcrqM3auhS
mQFuraiBxbfD7svRBDYLWZCEUkKhubGcCyAe3M6hSKgzv+SwO6TAXaYot1M5j/ZI
gECXA2tyKEInFK9dJJgR3fwjOqOBmTCBljAdBgNVHQ4EFgQU/GIWLL1zMisS3+Or
knZOFs7qgdswHwYDVR0jBBgwFoAU/GIWLL1zMisS3+OrknZOFs7qgdswDwYDVR0T
AQH/BAUwAwEB/zAhBgNVHREEGjAYghZmbGVldGFwaS5tb292ZXRyYXguY29tMBMG
A1UdJQQMMAoGCCsGAQUFBwMBMAsGA1UdDwQEAwICjDAKBggqhkjOPQQDAgOBjAAw
gYgCQgDB2t8wcLm+wFQA1uGYV124e8YHKKU5gHZZP1BNarKIDOkgg27nWvrTqEje
DKR57K6rGLQk04wfNG6AqqwKpYHCJQJCAWUIxP8V5afBj3OosLsH266CNrvU+0F8
g+nCMaavadk/8KExFVxQTkdyPLjZDQM2+VkReuzu9GteonN+zGa7cuVo
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

	teslaCredential.TestServerUri = "https://test.moovetrax.com"

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
