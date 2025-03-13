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
MIIDjjCCAxSgAwIBAgISA32xa9QDO/cRB+4+rHWs8kowMAoGCCqGSM49BAMDMDIx
CzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQDEwJF
NTAeFw0yNTAzMTIxNzI4MDhaFw0yNTA2MTAxNzI4MDdaMCExHzAdBgNVBAMTFmZs
ZWV0YXBpLm1vb3ZldHJheC5jb20wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAARk
QNrk+6qQce77IUcxGHHaYRDKxehyM+PfcE2B6JRfWH2mto+2aupJdnGBMtwIjCLi
cv5Zrh3czuEqcLE9znzVo4ICGTCCAhUwDgYDVR0PAQH/BAQDAgeAMB0GA1UdJQQW
MBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBT/
HjuyVIxM9KAiSGLX6U0v0lSV1TAfBgNVHSMEGDAWgBSfK1/PPCFPnQS37SssxMZw
i9LXDTBVBggrBgEFBQcBAQRJMEcwIQYIKwYBBQUHMAGGFWh0dHA6Ly9lNS5vLmxl
bmNyLm9yZzAiBggrBgEFBQcwAoYWaHR0cDovL2U1LmkubGVuY3Iub3JnLzAhBgNV
HREEGjAYghZmbGVldGFwaS5tb292ZXRyYXguY29tMBMGA1UdIAQMMAowCAYGZ4EM
AQIBMIIBBQYKKwYBBAHWeQIEAgSB9gSB8wDxAHYAzPsPaoVxCWX+lZtTzumyfCLp
hVwNl422qX5UwP5MDbAAAAGVi5vMIgAABAMARzBFAiAoJfPPUWfsyOSAb4AOCYz7
0bWfq320mBnGIxIslZpRXAIhAJkKxQ2wrJz8cWXMwFhtazXnzcWa/GVqFH0HB5RX
PwEKAHcA3oWB11AkfGvNy69WN8XngcZM5G7WF2OfjzSnJsnivTcAAAGVi5vMEAAA
BAMASDBGAiEA/PaPXco5xfLwZniTYTNlWYPQAPzmSFPzNi+KEJ04n+oCIQCoL6AI
gLAl+DU67nDh6j2DXra0uvpN/u9K4tkPaJhrWzAKBggqhkjOPQQDAwNoADBlAjBM
4XeUjlixhWazVQpGthuPD2ARTqK5gdETHpI2DWnoEktV4ech3uV+0R1IMQq2KA8C
MQDnSfP83DjlwJ7kHhkUksprUVtbIUtMnnbrHfo7dQltUHk+IXfv5QFnFWeenvqq
wBA=
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
