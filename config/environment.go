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
