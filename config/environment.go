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
MIICQTCCAaKgAwIBAgIUV1eDGL5tmRN+4hvJ3or/zzz9N9QwCgYIKoZIzj0EAwIw
ITEfMB0GA1UEAwwWZmxlZXRhcGkubW9vdmV0cmF4LmNvbTAeFw0yNTAzMTIyMjAx
MjBaFw0zNTAzMTAyMjAxMjBaMCExHzAdBgNVBAMMFmZsZWV0YXBpLm1vb3ZldHJh
eC5jb20wgZswEAYHKoZIzj0CAQYFK4EEACMDgYYABAHSdetrRPyocU1fkL9020xU
I/xB7OwsfR0X/vbkTzA4XpHYR2IlQR8NCBEp4j7Tk/RZyYUPCLzze7Gdm/uDlDhv
RgG+or8i6OWAHHZuyW/3SmNq7HoZmxT44SFIRnhSZARtHWNMi+ibn6fJUTdUQlP7
Sg82In8DNfHTh8cc++ILzjCZOaN1MHMwHQYDVR0OBBYEFI+ZLIzbfEHk4PuiuPqR
2ia+n2L+MB8GA1UdIwQYMBaAFI+ZLIzbfEHk4PuiuPqR2ia+n2L+MA8GA1UdEwEB
/wQFMAMBAf8wEwYDVR0lBAwwCgYIKwYBBQUHAwEwCwYDVR0PBAQDAgKMMAoGCCqG
SM49BAMCA4GMADCBiAJCAOSsSB3OkJzv5eng1iaghjq5ueOV7WEIIu6tWjBK6wDE
sO9T0EIkb+hZauXAyNB3jjJs/AaNfcFNXpVpdUasJAW1AkIBSHyYGchbXLUAxjsa
aHrnGdQNtmk2NnNr6qHn9N9+OlENc6DXhsdXpiYLcodLI/AlreI0IvKfrG91EiaY
DNKb0vo=
-----END CERTIFICATE-----` // tls-cert.pem
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
