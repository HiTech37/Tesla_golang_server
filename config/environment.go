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

	var tlsCertificate = `-----BEGIN PRIVATE KEY-----
MIHuAgEAMBAGByqGSM49AgEGBSuBBAAjBIHWMIHTAgEBBEIBtGo7VunUVBYBuwLC
pci7dyZ2hlR8J9FwxbC8Q5hnE5Wd6lAzhlQ6Ie423dJPqFr4L7wFybKyxtZ8Dw8H
hh1G2hyhgYkDgYYABACHy0AJ6hErMCQ4q9GS5q0uYYAtS88A/6iHQM4kIC+33GlN
lYIKSKk5OsTzE5azoPFMrGCuWAR/f4ou81pcABGMLAHfMnnqEu9e8RfKj6C2TW/w
0ABYD7ACRt4JMvVjzbC0PHTkyT+++Jt7E6cF2tdPoN2XG93T0GItFBRTzSo2wKzJ
VQ==
-----END PRIVATE KEY-----`

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

// func FormatCertFile() string {
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
