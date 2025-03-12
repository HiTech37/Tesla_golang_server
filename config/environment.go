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
MIICJzCCAYigAwIBAgIUV5d1aGw0/8U6Tm4HjHIsE1BRlrcwCgYIKoZIzj0EAwIw
FDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI1MDMxMjIxMDIzMVoXDTM1MDMxMDIx
MDIzMVowFDESMBAGA1UEAwwJbG9jYWxob3N0MIGbMBAGByqGSM49AgEGBSuBBAAj
A4GGAAQAT9or41S2V5gSWe2bBr3w/iPnQvOZ5dXOncaZONVzhotMWEnxxQKgOYwR
Bia3uBj+tNQjeNBFsF8CnAXYgqFCWEkBLZxLADIsONEwk32PnicUcgRsGCZ09HH/
LTKwqwJzndzhYu1AbG7vaS+oep3YqNDHe8H5kDP89FrATM0TYZUKnHOjdTBzMB0G
A1UdDgQWBBQDiZ7HxUzp1FAhpG8uO3GPaAasrzAfBgNVHSMEGDAWgBQDiZ7HxUzp
1FAhpG8uO3GPaAasrzAPBgNVHRMBAf8EBTADAQH/MBMGA1UdJQQMMAoGCCsGAQUF
BwMBMAsGA1UdDwQEAwICjDAKBggqhkjOPQQDAgOBjAAwgYgCQgDE7VN03aQv/UBM
2Z3JF/3VEMlWbmO6OPlH6IvqtkC6+WF/m6tb89ES/s7plmc0S1K6fx9I11kxkv3x
8RlXS658+gJCAfyoVgiiBJqKoVFBMEz4yqhthI+jTIGcsmZ2k4SmdTuC5MqMBG32
lbjMiKiRso9D6bwXnZgj23eesDyesbjgjqBx
-----END CERTIFICATE-----` // cert.pem
	teslaCredential.CallbackUri = "https://fleetapi.moovetrax.com/device_signup"
	teslaCredential.ClientID = "60d97918-9b6b-4c92-88e3-ff9e9403239f"
	teslaCredential.SecretKey = "ta-secret.8p8Jz&Y^%n9FCCeE"
	teslaCredential.RootDomain = "https://moovetrax.com"
	teslaCredential.DataScope = "openid%20vehicle_device_data%20vehicle_location%20offline_access%20vehicle_cmds%20vehicle_charging_cmds"
	teslaCredential.Certificate = certificate
	teslaCredential.ServerDomain = "fleetapi.moovetrax.com"
	teslaCredential.Port = 8443
	teslaCredential.ProxyUri = "https://vehicle-command-tesla_http_proxy-1:4443"
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
