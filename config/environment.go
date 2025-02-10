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

	var environment = "test" // "local", "test", "prod"

	if environment == "local" {
		teslaCredential.CallbackUri = "http://localhost:3000/tesla_signup"
	} else if environment == "test" {
		teslaCredential.CallbackUri = "https://test.moovetrax.com/tesla_signup"

	} else if environment == "prod" {
		teslaCredential.CallbackUri = "https://moovetrax.com/tesla_signup"
	}

	var certificate = `-----BEGIN CERTIFICATE-----
MIIDjDCCAxKgAwIBAgISA3wwTpiRnSbENZ9NQ1LD2xVNMAoGCCqGSM49BAMDMDIx
CzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQDEwJF
NTAeFw0yNTAyMTAwMzI5NDZaFw0yNTA1MTEwMzI5NDVaMCExHzAdBgNVBAMTFnRl
c2xhYXBpLm1vb3ZldHJheC5jb20wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAT/
OTB8jEE8W9oVkjTe7Guv/Otnc/tpTNxZVWKehMYIwqEGJ96ZjaR191YkfhEgMDKg
MFlqejLvZd/Ag34LV+ULo4ICFzCCAhMwDgYDVR0PAQH/BAQDAgeAMB0GA1UdJQQW
MBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBTI
McEBHuZMlaItfuIWnGJ8JPujSDAfBgNVHSMEGDAWgBSfK1/PPCFPnQS37SssxMZw
i9LXDTBVBggrBgEFBQcBAQRJMEcwIQYIKwYBBQUHMAGGFWh0dHA6Ly9lNS5vLmxl
bmNyLm9yZzAiBggrBgEFBQcwAoYWaHR0cDovL2U1LmkubGVuY3Iub3JnLzAhBgNV
HREEGjAYghZ0ZXNsYWFwaS5tb292ZXRyYXguY29tMBMGA1UdIAQMMAowCAYGZ4EM
AQIBMIIBAwYKKwYBBAHWeQIEAgSB9ASB8QDvAHUAouMK5EXvva2bfjjtR2d3U9eC
W4SU1yteGyzEuVCkR+cAAAGU7h16uAAABAMARjBEAiAiaMq8VDB6aLhOKw0aN8zD
nkPAy/mc0WaltjyWjptLoQIgCIhu64wA70eNl/qSjt08nBnvTmX+iEcuIVz2Jnbh
r1EAdgDPEVbu1S58r/OHW9lpLpvpGnFnSrAX7KwB0lt3zsw7CAAAAZTuHXrzAAAE
AwBHMEUCIQC613iEbznMx5zoARHIfBGnI5hDQSucztC6XpJxnXtBmwIgVGudXPrW
2DFz1aGVcdp3ytt/XYCZNi7phjvapaaN3JMwCgYIKoZIzj0EAwMDaAAwZQIwRL0O
S8EuCN3/wlk/MWYE/rQNEecjhiKRC3acPm3e22fert0ChUR1AFK1f9U6YRRyAjEA
6igmXHRnoCf+YwArdfrz+ltqMh1nLRkLbvTKsqtVtJpi+8rzqNPaVccNF/CtteKV
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIEVzCCAj+gAwIBAgIRAIOPbGPOsTmMYgZigxXJ/d4wDQYJKoZIhvcNAQELBQAw
TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh
cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMjQwMzEzMDAwMDAw
WhcNMjcwMzEyMjM1OTU5WjAyMQswCQYDVQQGEwJVUzEWMBQGA1UEChMNTGV0J3Mg
RW5jcnlwdDELMAkGA1UEAxMCRTUwdjAQBgcqhkjOPQIBBgUrgQQAIgNiAAQNCzqK
a2GOtu/cX1jnxkJFVKtj9mZhSAouWXW0gQI3ULc/FnncmOyhKJdyIBwsz9V8UiBO
VHhbhBRrwJCuhezAUUE8Wod/Bk3U/mDR+mwt4X2VEIiiCFQPmRpM5uoKrNijgfgw
gfUwDgYDVR0PAQH/BAQDAgGGMB0GA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcD
ATASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1UdDgQWBBSfK1/PPCFPnQS37SssxMZw
i9LXDTAfBgNVHSMEGDAWgBR5tFnme7bl5AFzgAiIyBpY9umbbjAyBggrBgEFBQcB
AQQmMCQwIgYIKwYBBQUHMAKGFmh0dHA6Ly94MS5pLmxlbmNyLm9yZy8wEwYDVR0g
BAwwCjAIBgZngQwBAgEwJwYDVR0fBCAwHjAcoBqgGIYWaHR0cDovL3gxLmMubGVu
Y3Iub3JnLzANBgkqhkiG9w0BAQsFAAOCAgEAH3KdNEVCQdqk0LKyuNImTKdRJY1C
2uw2SJajuhqkyGPY8C+zzsufZ+mgnhnq1A2KVQOSykOEnUbx1cy637rBAihx97r+
bcwbZM6sTDIaEriR/PLk6LKs9Be0uoVxgOKDcpG9svD33J+G9Lcfv1K9luDmSTgG
6XNFIN5vfI5gs/lMPyojEMdIzK9blcl2/1vKxO8WGCcjvsQ1nJ/Pwt8LQZBfOFyV
XP8ubAp/au3dc4EKWG9MO5zcx1qT9+NXRGdVWxGvmBFRAajciMfXME1ZuGmk3/GO
koAM7ZkjZmleyokP1LGzmfJcUd9s7eeu1/9/eg5XlXd/55GtYjAM+C4DG5i7eaNq
cm2F+yxYIPt6cbbtYVNJCGfHWqHEQ4FYStUyFnv8sjyqU8ypgZaNJ9aVcWSICLOI
E1/Qv/7oKsnZCWJ926wU6RqG1OYPGOi1zuABhLw61cuPVDT28nQS/e6z95cJXq0e
K1BcaJ6fJZsmbjRgD5p3mvEf5vdQM7MCEvU0tHbsx2I5mHHJoABHb8KVBgWp/lcX
GWiWaeOyB7RP+OfDtvi2OsapxXiV7vNVs7fMlrRjY1joKaqmmycnBvAq14AEbtyL
sVfOS66B8apkeFX2NY4XPEYV4ZSCe8VHPrdrERk2wILG3T/EGmSIkCYVUMSnjmJd
VQD9F6Na/+zmXCc=
-----END CERTIFICATE-----`

	var tlsCertificate = `-----BEGIN CERTIFICATE-----
MIICZjCCAcegAwIBAgIUEjutwFXrDSFfT9rGPM1JsUGRWBAwCgYIKoZIzj0EAwIw
ITEfMB0GA1UEAwwWdDNzbGFhcGkubW9vdmV0cmF4LmNvbTAeFw0yNTAxMjMxODQ1
NDZaFw0zNTAxMjExODQ1NDZaMCExHzAdBgNVBAMMFnQzc2xhYXBpLm1vb3ZldHJh
eC5jb20wgZswEAYHKoZIzj0CAQYFK4EEACMDgYYABACHy0AJ6hErMCQ4q9GS5q0u
YYAtS88A/6iHQM4kIC+33GlNlYIKSKk5OsTzE5azoPFMrGCuWAR/f4ou81pcABGM
LAHfMnnqEu9e8RfKj6C2TW/w0ABYD7ACRt4JMvVjzbC0PHTkyT+++Jt7E6cF2tdP
oN2XG93T0GItFBRTzSo2wKzJVaOBmTCBljAdBgNVHQ4EFgQU0Ux6uBBaZKMjQK/9
S/FVwFIu0+IwHwYDVR0jBBgwFoAU0Ux6uBBaZKMjQK/9S/FVwFIu0+IwDwYDVR0T
AQH/BAUwAwEB/zAhBgNVHREEGjAYghZ0M3NsYWFwaS5tb292ZXRyYXguY29tMBMG
A1UdJQQMMAoGCCsGAQUFBwMBMAsGA1UdDwQEAwICjDAKBggqhkjOPQQDAgOBjAAw
gYgCQgGH92HrWa16wd8Dk8eMK/IZiHlAQn5MXGQP/Bua1+SaCOe6qUg2VviPtDIn
9+rMeA5MO/iDWorhYRpd0efMEry1ZwJCAMh/67W5/2XkRjvi7s+LtD11BXtJhT+C
yLNK9r/++lSTTpZPEiEbBM+Lnyvh0FfQNZ8qImt3kRKcwqYyramlHXWU
-----END CERTIFICATE-----`

	var clientCert = `-----BEGIN CERTIFICATE-----
MIICZjCCAcegAwIBAgIUNnV3TC4SvTTm4IqH/tPtq7K75ucwCgYIKoZIzj0EAwIw
ITEfMB0GA1UEAwwWdDNzbGFhcGkubW9vdmV0cmF4LmNvbTAeFw0yNTAxMjMyMzIx
MzhaFw0zNTAxMjEyMzIxMzhaMCExHzAdBgNVBAMMFnQzc2xhYXBpLm1vb3ZldHJh
eC5jb20wgZswEAYHKoZIzj0CAQYFK4EEACMDgYYABAG1cIOI/+9deCugmK9niPuG
ie8+7WOvd7LkIPQH4mnXy1IMXq55uJgf5pCdSi/vuG6ghs+vY9WkRL0HzmOApN4y
cABQw0372WNiHna+NrJgw7P1Qj8tXW1prdRY8eiGrd3A/BHyFM6C/ql/KzBCkSHc
g2GhxOkvu6tIqAwxTddCjvVbq6OBmTCBljAdBgNVHQ4EFgQU55WltTHc8PiiyYHw
7oUArgxIihMwHwYDVR0jBBgwFoAU55WltTHc8PiiyYHw7oUArgxIihMwDwYDVR0T
AQH/BAUwAwEB/zAhBgNVHREEGjAYghZ0M3NsYWFwaS5tb292ZXRyYXguY29tMBMG
A1UdJQQMMAoGCCsGAQUFBwMBMAsGA1UdDwQEAwICjDAKBggqhkjOPQQDAgOBjAAw
gYgCQgFWqLABpuCqMN9k/iIJtO5nSl67Xwa3YPcgaeUj06Kf9GxHO4D08lbCwoDb
3G14cWWumRqzpE6Kvkjpo6ziULMlnwJCAPy+IWGk5syF77MWz4UbEznVQZU2AWtX
M1PItUSAXh5yYEgQdBgYag7pdkesyoL9SmfayRR9ZqQXjJe14sxpha94
-----END CERTIFICATE-----`

	var clientKey = `-----BEGIN PRIVATE KEY-----
MIHuAgEAMBAGByqGSM49AgEGBSuBBAAjBIHWMIHTAgEBBEIBknLkXISDbyxWPceT
mlcZtpTAU6LRWmnCw/5y0kBzpRv+CAPOV+nisPhKMUXgMmlY0RvuW1RD+SCfLVOP
rZYBe9ChgYkDgYYABAG1cIOI/+9deCugmK9niPuGie8+7WOvd7LkIPQH4mnXy1IM
Xq55uJgf5pCdSi/vuG6ghs+vY9WkRL0HzmOApN4ycABQw0372WNiHna+NrJgw7P1
Qj8tXW1prdRY8eiGrd3A/BHyFM6C/ql/KzBCkSHcg2GhxOkvu6tIqAwxTddCjvVb
qw==
-----END PRIVATE KEY-----`

	teslaCredential.ClientID = "69e55814-1679-46d3-a3b6-ac713f77f287"
	teslaCredential.SecretKey = "ta-secret.TjmkFpMgD_pXgdBA"
	teslaCredential.RootDomain = "https://moovetrax.com"
	teslaCredential.DataScope = "openid vehicle_device_data vehicle_location offline_access vehicle_cmds vehicle_charging_cmds"
	teslaCredential.Certificate = certificate
	teslaCredential.TlsCertificate = tlsCertificate
	teslaCredential.ServerDomain = "teslaapi.moovetrax.com"
	teslaCredential.Port = 8443
	teslaCredential.ProxyUri = "https://teslaapi.moovetrax.com:4443"
	teslaCredential.ClientCert = clientCert
	teslaCredential.ClientKey = clientKey

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
