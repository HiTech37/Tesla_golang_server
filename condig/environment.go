package config

type TeslaCredential struct {
	ClientID     string
	SecretKey    string
	CallbackUri  string
	DataScope    string
	RootDomain   string
	ServerDomain string
	Certificate  string
	ProxyUri     string
	Port         int
}

func GetTeslaCredential() *TeslaCredential {
	var teslaCredential TeslaCredential

	var certificate = ""

	var environment = "local" // "local", "test", "prod"

	if environment == "local" {
		teslaCredential.CallbackUri = "http://localhost:3000/tesla_signup"
	} else if environment == "test" {
		teslaCredential.CallbackUri = "https://test.moovetrax.com/tesla_signup"
	} else if environment == "prod" {
		teslaCredential.CallbackUri = "https://moovetrax.com/tesla_signup"
	}

	teslaCredential.ClientID = "69e55814-1679-46d3-a3b6-ac713f77f287"
	teslaCredential.SecretKey = "ta-secret.TjmkFpMgD_pXgdBA"
	teslaCredential.RootDomain = "https://moovetrax.com"
	teslaCredential.DataScope = "openid vehicle_device_data vehicle_location offline_access vehicle_cmds vehicle_charging_cmds"
	teslaCredential.Certificate = certificate
	teslaCredential.ServerDomain = "t3slaapi.moovetrax.com"
	teslaCredential.Port = 4443
	teslaCredential.ProxyUri = "https://localhost:4443"

	return &teslaCredential
}
