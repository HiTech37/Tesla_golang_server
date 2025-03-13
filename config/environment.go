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
MIICgjCCAeOgAwIBAgIUM9jdAxMdRLkKvZ4MfciDE0Ldhs4wCgYIKoZIzj0EAwIw
ITEfMB0GA1UEAwwWZmxlZXRhcGkubW9vdmV0cmF4LmNvbTAeFw0yNTAzMTMyMTMx
MzNaFw0zNTAzMTEyMTMxMzNaMCExHzAdBgNVBAMMFmZsZWV0YXBpLm1vb3ZldHJh
eC5jb20wgZswEAYHKoZIzj0CAQYFK4EEACMDgYYABAEc1y1Wb2LSoSx7PIxIAsNn
amsWH+rJmdo8uXh9VHukZQlO77bu32qVB/9lbd97kiJoTFYv3cmchoiJdYrskIIK
/gEmPUwPvSaurMZfxQ8uBst7Ex3Lq9TuQ+lG0w44DSvxlXz+iiogAYcUT+X4hsAV
yz+xGzGqNQuXNKNKV3zw5S5eF6OBtTCBsjAdBgNVHQ4EFgQUQpUvsGxUmWG3dSAS
xfWjH/HIpY4wHwYDVR0jBBgwFoAUQpUvsGxUmWG3dSASxfWjH/HIpY4wDwYDVR0T
AQH/BAUwAwEB/zATBgNVHSUEDDAKBggrBgEFBQcDATALBgNVHQ8EBAMCAowwPQYD
VR0RBDYwNIIWZmxlZXRhcGkubW9vdmV0cmF4LmNvbYIad3d3LmZsZWV0YXBpLm1v
b3ZldHJheC5jb20wCgYIKoZIzj0EAwIDgYwAMIGIAkIBfiRuZQT0bh5Ql3RBnx0/
5btCSClL11yJ6ZCxqlDRXRZbLbZnUwH8Y6KoFya3g80YnIPcO2WNLjluTP0fHsw/
JCgCQgE6pf2XP6elzutNrQ90aDYo+Tp8pGeofmtebhZwRrHISpFgE5bWj8QTZ38z
+PDP9hBP0rHGoUohH+vEosS0T224YA==
-----END CERTIFICATE-----` // cert.pem
	teslaCredential.CallbackUri = "https://fleetapi.moovetrax.com/device_signup"
	teslaCredential.ClientID = "60d97918-9b6b-4c92-88e3-ff9e9403239f"
	teslaCredential.SecretKey = "ta-secret.8p8Jz&Y^%n9FCCeE"
	teslaCredential.RootDomain = "https://moovetrax.com"
	teslaCredential.DataScope = "openid%20vehicle_device_data%20vehicle_location%20offline_access%20vehicle_cmds%20vehicle_charging_cmds"
	teslaCredential.Certificate = certificate
	teslaCredential.ServerDomain = "fleetapi.moovetrax.com"
	teslaCredential.Port = 8443
	teslaCredential.ProxyUri = "https://localhost:4443"
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
