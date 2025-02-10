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
MIIDjDCCAxOgAwIBAgISAw/1GAcf/EII5oESrdZNlkpsMAoGCCqGSM49BAMDMDIx
CzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQDEwJF
NjAeFw0yNTAyMTAwODA3MzNaFw0yNTA1MTEwODA3MzJaMCExHzAdBgNVBAMTFnRl
c2xhYXBpLm1vb3ZldHJheC5jb20wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAARx
zWNAs3qW1RwQuaeFfDz3+9KFuTWiHr9odbWE0Eb78v5rdjpQa/1YUzAoRAOfn8Tt
89wL51iR3TYQNsBeLbS5o4ICGDCCAhQwDgYDVR0PAQH/BAQDAgeAMB0GA1UdJQQW
MBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBQ8
LiKdD4kxWnC/f8WI5lgEWb6+xTAfBgNVHSMEGDAWgBSTJ0aYA6lRaI6Y1sRCSNsj
v1iU0jBVBggrBgEFBQcBAQRJMEcwIQYIKwYBBQUHMAGGFWh0dHA6Ly9lNi5vLmxl
bmNyLm9yZzAiBggrBgEFBQcwAoYWaHR0cDovL2U2LmkubGVuY3Iub3JnLzAhBgNV
HREEGjAYghZ0ZXNsYWFwaS5tb292ZXRyYXguY29tMBMGA1UdIAQMMAowCAYGZ4EM
AQIBMIIBBAYKKwYBBAHWeQIEAgSB9QSB8gDwAHYAzPsPaoVxCWX+lZtTzumyfCLp
hVwNl422qX5UwP5MDbAAAAGU7xvKBQAABAMARzBFAiB83/lFsUrj0XVshL6v3hts
d3WY6ROHiGrKCpSHgFS2QwIhAJBllDTEjfNwRPY6EXtFrE/MYj0UKNHIgkuU9zG7
MVIeAHYAzxFW7tUufK/zh1vZaS6b6RpxZ0qwF+ysAdJbd87MOwgAAAGU7xvKMQAA
BAMARzBFAiEAxAedHBkyFXWkE3ul65TtlZ9FiBVTATjqM+5171C712oCIHoRkmWz
YYJRZaxDSyplFjzkOyEm3HWAnapln+L3DGQIMAoGCCqGSM49BAMDA2cAMGQCMFfi
lAV1fkAXvlhxrZxPHLRwJaCSnBQsSNtgbD5JoTRdQ3GiAGkkyMxioRdWaw6kxwIw
O9vO/9HJAwy7rJwvBOQMNPM7STTxE+BpRwBmigeTBj29vsUGJ+qeODv+CTScocq3
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
MIIDjDCCAxOgAwIBAgISAw/1GAcf/EII5oESrdZNlkpsMAoGCCqGSM49BAMDMDIx
CzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQDEwJF
NjAeFw0yNTAyMTAwODA3MzNaFw0yNTA1MTEwODA3MzJaMCExHzAdBgNVBAMTFnRl
c2xhYXBpLm1vb3ZldHJheC5jb20wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAARx
zWNAs3qW1RwQuaeFfDz3+9KFuTWiHr9odbWE0Eb78v5rdjpQa/1YUzAoRAOfn8Tt
89wL51iR3TYQNsBeLbS5o4ICGDCCAhQwDgYDVR0PAQH/BAQDAgeAMB0GA1UdJQQW
MBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBQ8
LiKdD4kxWnC/f8WI5lgEWb6+xTAfBgNVHSMEGDAWgBSTJ0aYA6lRaI6Y1sRCSNsj
v1iU0jBVBggrBgEFBQcBAQRJMEcwIQYIKwYBBQUHMAGGFWh0dHA6Ly9lNi5vLmxl
bmNyLm9yZzAiBggrBgEFBQcwAoYWaHR0cDovL2U2LmkubGVuY3Iub3JnLzAhBgNV
HREEGjAYghZ0ZXNsYWFwaS5tb292ZXRyYXguY29tMBMGA1UdIAQMMAowCAYGZ4EM
AQIBMIIBBAYKKwYBBAHWeQIEAgSB9QSB8gDwAHYAzPsPaoVxCWX+lZtTzumyfCLp
hVwNl422qX5UwP5MDbAAAAGU7xvKBQAABAMARzBFAiB83/lFsUrj0XVshL6v3hts
d3WY6ROHiGrKCpSHgFS2QwIhAJBllDTEjfNwRPY6EXtFrE/MYj0UKNHIgkuU9zG7
MVIeAHYAzxFW7tUufK/zh1vZaS6b6RpxZ0qwF+ysAdJbd87MOwgAAAGU7xvKMQAA
BAMARzBFAiEAxAedHBkyFXWkE3ul65TtlZ9FiBVTATjqM+5171C712oCIHoRkmWz
YYJRZaxDSyplFjzkOyEm3HWAnapln+L3DGQIMAoGCCqGSM49BAMDA2cAMGQCMFfi
lAV1fkAXvlhxrZxPHLRwJaCSnBQsSNtgbD5JoTRdQ3GiAGkkyMxioRdWaw6kxwIw
O9vO/9HJAwy7rJwvBOQMNPM7STTxE+BpRwBmigeTBj29vsUGJ+qeODv+CTScocq3
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIEVzCCAj+gAwIBAgIRALBXPpFzlydw27SHyzpFKzgwDQYJKoZIhvcNAQELBQAw
TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh
cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMjQwMzEzMDAwMDAw
WhcNMjcwMzEyMjM1OTU5WjAyMQswCQYDVQQGEwJVUzEWMBQGA1UEChMNTGV0J3Mg
RW5jcnlwdDELMAkGA1UEAxMCRTYwdjAQBgcqhkjOPQIBBgUrgQQAIgNiAATZ8Z5G
h/ghcWCoJuuj+rnq2h25EqfUJtlRFLFhfHWWvyILOR/VvtEKRqotPEoJhC6+QJVV
6RlAN2Z17TJOdwRJ+HB7wxjnzvdxEP6sdNgA1O1tHHMWMxCcOrLqbGL0vbijgfgw
gfUwDgYDVR0PAQH/BAQDAgGGMB0GA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcD
ATASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1UdDgQWBBSTJ0aYA6lRaI6Y1sRCSNsj
v1iU0jAfBgNVHSMEGDAWgBR5tFnme7bl5AFzgAiIyBpY9umbbjAyBggrBgEFBQcB
AQQmMCQwIgYIKwYBBQUHMAKGFmh0dHA6Ly94MS5pLmxlbmNyLm9yZy8wEwYDVR0g
BAwwCjAIBgZngQwBAgEwJwYDVR0fBCAwHjAcoBqgGIYWaHR0cDovL3gxLmMubGVu
Y3Iub3JnLzANBgkqhkiG9w0BAQsFAAOCAgEAfYt7SiA1sgWGCIpunk46r4AExIRc
MxkKgUhNlrrv1B21hOaXN/5miE+LOTbrcmU/M9yvC6MVY730GNFoL8IhJ8j8vrOL
pMY22OP6baS1k9YMrtDTlwJHoGby04ThTUeBDksS9RiuHvicZqBedQdIF65pZuhp
eDcGBcLiYasQr/EO5gxxtLyTmgsHSOVSBcFOn9lgv7LECPq9i7mfH3mpxgrRKSxH
pOoZ0KXMcB+hHuvlklHntvcI0mMMQ0mhYj6qtMFStkF1RpCG3IPdIwpVCQqu8GV7
s8ubknRzs+3C/Bm19RFOoiPpDkwvyNfvmQ14XkyqqKK5oZ8zhD32kFRQkxa8uZSu
h4aTImFxknu39waBxIRXE4jKxlAmQc4QjFZoq1KmQqQg0J/1JF8RlFvJas1VcjLv
YlvUB2t6npO6oQjB3l+PNf0DpQH7iUx3Wz5AjQCi6L25FjyE06q6BZ/QlmtYdl/8
ZYao4SRqPEs/6cAiF+Qf5zg2UkaWtDphl1LKMuTNLotvsX99HP69V2faNyegodQ0
LyTApr/vT01YPE46vNsDLgK+4cL6TrzC/a4WcmF5SRJ938zrv/duJHLXQIku5v0+
EwOy59Hdm0PT/Er/84dDV0CSjdR/2XuZM3kpysSKLgD1cKiDA+IRguODCxfO9cyY
Ig46v9mFmBvyH04=
-----END CERTIFICATE-----`

	var clientKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIA2/TNhjaS7dc7MiW3hZjoqf8hs6jiCYGQbYRPTvCAenoAoGCCqGSM49
AwEHoUQDQgAEcc1jQLN6ltUcELmnhXw89/vShbk1oh6/aHW1hNBG+/L+a3Y6UGv9
WFMwKEQDn5/E7fPcC+dYkd02EDbAXi20uQ==
-----END EC PRIVATE KEY-----`

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
