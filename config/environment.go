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
