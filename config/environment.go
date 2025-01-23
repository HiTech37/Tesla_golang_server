package config

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
MIICPzCCAaCgAwIBAgIUO5pfC+X9zzBPND4PCo0mjSMCKcswCgYIKoZIzj0EAwIw
FDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI1MDEyMzA4MzA0OFoXDTM1MDEyMTA4
MzA0OFowFDESMBAGA1UEAwwJbG9jYWxob3N0MIGbMBAGByqGSM49AgEGBSuBBAAj
A4GGAAQBIdT5I4cbnzH8i8NNqfZD702Bk4tBe8BtZm/NleNVAt0btWgWqQQRn+Gh
z+NN5m5Mg9dGZIg3Mo2wnfqu/9683X0BcG5cWst4OtYGvVGQU7/xya4mx0nsXJUN
KKvu7LThlWHk6VdIoG4kkt0bwtgw4LMR171BlMcm69hAlT7y+ldmMDqjgYwwgYkw
HQYDVR0OBBYEFDVTUKjEeVY6ZcPGpRJAZMLfFE+lMB8GA1UdIwQYMBaAFDVTUKjE
eVY6ZcPGpRJAZMLfFE+lMA8GA1UdEwEB/wQFMAMBAf8wFAYDVR0RBA0wC4IJbG9j
YWxob3N0MBMGA1UdJQQMMAoGCCsGAQUFBwMBMAsGA1UdDwQEAwICjDAKBggqhkjO
PQQDAgOBjAAwgYgCQgGS1a+NgvaBTgbRuzdphegwj8AJtNkXCfYXJsMTbBF5Sul0
GXdxMAkcTIRzfunTzOFijd9XRRw2mn0FOtf0mbHUmAJCAQtOFXEkj07jw8cbYsU5
JqjhPIzL62dVmVKZ4DibvtP5ZnVjb8fSRfcBHuSbIR7vPFYhRmCZsc8CHHwyTqz7
Tf+Y
-----END CERTIFICATE-----`

	var tlsCertificate = `-----BEGIN CERTIFICATE-----
MIIDjjCCAxOgAwIBAgISBJRSih/BQ1EP/CEDb7GKJSKFMAoGCCqGSM49BAMDMDIx
CzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQDEwJF
NTAeFw0yNTAxMjIxNzQyMTBaFw0yNTA0MjIxNzQyMDlaMCExHzAdBgNVBAMTFnQz
c2xhYXBpLm1vb3ZldHJheC5jb20wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATz
CFM6huzLJCPbkH6pSjSJy+Wxt3G2UL4AfZbGjuBbMz0pmIzr7Zi4dhkMtJ6JYynS
boKJejMbccabckMo/iaNo4ICGDCCAhQwDgYDVR0PAQH/BAQDAgeAMB0GA1UdJQQW
MBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBRU
Mn0U5jowA3unqZPZWRS4EAkxmzAfBgNVHSMEGDAWgBSfK1/PPCFPnQS37SssxMZw
i9LXDTBVBggrBgEFBQcBAQRJMEcwIQYIKwYBBQUHMAGGFWh0dHA6Ly9lNS5vLmxl
bmNyLm9yZzAiBggrBgEFBQcwAoYWaHR0cDovL2U1LmkubGVuY3Iub3JnLzAhBgNV
HREEGjAYghZ0M3NsYWFwaS5tb292ZXRyYXguY29tMBMGA1UdIAQMMAowCAYGZ4EM
AQIBMIIBBAYKKwYBBAHWeQIEAgSB9QSB8gDwAHcA5tIxY0B3jMEQQQbXcbnOwdJA
9paEhvu6hzId/R43jlAAAAGUj1ELBwAABAMASDBGAiEAmKpZfKj6KX8Fs+u+MKtn
3hN5Xhc03IClG/sYA7WnfhMCIQCOMDCFwBvOZkiTkBTNAon/0hoZGpbiJyWxzfen
QwETWQB1AKLjCuRF772tm3447Udnd1PXgluElNcrXhssxLlQpEfnAAABlI9RCwoA
AAQDAEYwRAIgFRnNa2CHidWiJgF4F6pDdHQQmIGBZTk9zVNDIS4D+coCIBTLLiBM
ZFjZHYbTZZVYdvAxG9SG2psObTU8W2YHH3BJMAoGCCqGSM49BAMDA2kAMGYCMQDG
/+bPlrwRN/Z5lx0Ye7HgMIORf7wDZBpNLT22VeQqZu7Nop3pUmuD9D/kbDWoFHUC
MQCMRJ6ZFKxDGAMAeCvZokcAAhgv3YddM35NNalXhAvDyCrbQV43mPEBBvyryR/M
MrA=
-----END CERTIFICATE-----`

	teslaCredential.ClientID = "69e55814-1679-46d3-a3b6-ac713f77f287"
	teslaCredential.SecretKey = "ta-secret.TjmkFpMgD_pXgdBA"
	teslaCredential.RootDomain = "https://moovetrax.com"
	teslaCredential.DataScope = "openid vehicle_device_data vehicle_location offline_access vehicle_cmds vehicle_charging_cmds"
	teslaCredential.Certificate = certificate
	teslaCredential.TlsCertificate = tlsCertificate
	teslaCredential.ServerDomain = "t3slaapi.moovetrax.com"
	teslaCredential.Port = 8443
	teslaCredential.ProxyUri = "https://localhost:4443"

	return &teslaCredential
}

// func ReadCertFile() string {
// 	filepath, err := filepath.Abs("./config/cert.pem")

// 	if err != nil {
// 		return ""
// 	}
// 	content, err := ioutil.ReadFile(filepath)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return ""
// 	}

// 	// formattedCertificate := strings.ReplaceAll(string(content), "\n", "\\n")

// 	// // Print in desired format
// 	// result := fmt.Sprintf("\"%s\"", formattedCertificate)
// 	return string(content)
// }
