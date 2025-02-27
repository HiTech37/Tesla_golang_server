package controller

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	config "tesla_server/config"
	"tesla_server/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestConnectParams struct {
	Vins         []string `json:"vins"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
}
type FieldConfig struct {
	ResendIntervalSeconds int `json:"resend_interval_seconds"`
	MinimumDelta          int `json:"minimum_delta"`
	IntervalSeconds       int `json:"interval_seconds"`
}

type Config struct {
	PreferTyped bool                   `json:"prefer_typed"`
	Port        int                    `json:"port"`
	Exp         int64                  `json:"exp"`
	AlertTypes  []string               `json:"alert_types"`
	Fields      map[string]FieldConfig `json:"fields"`
	CA          string                 `json:"ca"`
	Hostname    string                 `json:"hostname"`
}

type TelemetryRequest struct {
	Config Config   `json:"config"`
	Vins   []string `json:"vins"`
}

type VehicleInfo struct {
	Vin         string `json:"vin"`
	DeviceName  string `json:"display_name"`
	VehicleID   string `json:"id_s"`
	CarType     string `json:"model"`
	Color       string `json:"color"`
	ShareAbi    bool   `json:"abi"`
	ShareTintAi bool   `json:"tint_ai"`
}

type DeviceInfoParams struct {
	Email        string          `json:"email"`
	DeviceList   []VehicleInfo   `json:"deviceList"`
	AccessToken  string          `json:"accessToken"`
	RefreshToken string          `json:"refreshToken"`
	ShareStatus  map[string]bool `json:"checkStatus"`
}

type VehicleInfoParams struct {
	Response struct {
		Color         string `json:"color"`
		VehicleID     int    `json:"vehicle_id"`
		VehicleConfig struct {
			CarType string `json:"car_type"`
		} `json:"vehicle_config"`
		VehicleState struct {
			CarVersion string `json:"car_version"`
		} `json:"vehicle_state"`
	} `json:"response"`
}

type Payload struct {
	Email        string          `json:"email"`
	DeviceList   []VehicleInfo   `json:"deviceList"`
	AccessToken  string          `json:"accessToken"`
	RefreshToken string          `json:"refreshToken"`
	ShareInfo    map[string]bool `json:"shareInfo"`
}

type SkippedVehicles struct {
	MissingKey          []string `json:"missing_key"`
	UnsupportedHardware []string `json:"unsupported_hardware"`
	UnsupportedFirmware []string `json:"unsupported_firmware"`
	MaxConfigs          []string `json:"max_configs"`
}

type ResponseData struct {
	UpdatedVehicles int              `json:"updated_vehicles"`
	SkippedVehicles *SkippedVehicles `json:"skipped_vehicles,omitempty"`
}

type Root struct {
	Response ResponseData `json:"response"`
}

func ConnectDeviceforTest(c *gin.Context) {
	var requestParams RequestConnectParams
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	base := config.GetTeslaCredential().ProxyUri
	path := "/api/1/vehicles/fleet_telemetry_config"
	url := fmt.Sprintf("%s%s", base, path)

	fieldToStream1 := "Location"
	fieldToStream2 := "BatteryLevel"
	fieldToStream3 := "VehicleSpeed"
	fieldToStream4 := "Odometer"

	telemetryData := TelemetryRequest{
		Config: Config{
			PreferTyped: true,
			Port:        8443,
			Exp:         1770670000,
			AlertTypes:  []string{"service"},
			Fields: map[string]FieldConfig{
				fieldToStream1: {
					ResendIntervalSeconds: 3600,
					MinimumDelta:          1,
					IntervalSeconds:       60,
				},
				fieldToStream2: {
					ResendIntervalSeconds: 3600,
					MinimumDelta:          1,
					IntervalSeconds:       60,
				},
				fieldToStream3: {
					ResendIntervalSeconds: 3600,
					MinimumDelta:          1,
					IntervalSeconds:       60,
				},
				fieldToStream4: {
					ResendIntervalSeconds: 3600,
					MinimumDelta:          1,
					IntervalSeconds:       60,
				},
			},
			CA:       config.GetTeslaCredential().Certificate,
			Hostname: config.GetTeslaCredential().ServerDomain,
		},
		Vins: requestParams.Vins,
	}

	payload, err := json.Marshal(telemetryData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Failed to marshal JSON data:",
			"error": err,
		})
		return
	}

	certPEM := config.GetTeslaCredential().Certificate

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(certPEM)); !ok {
		log.Fatal("Failed to append certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		ServerName: config.GetTeslaCredential().ServerDomain,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error creating request:",
			"error": err,
		})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", requestParams.AccessToken))

	// client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error making request:",
			"error": err,
		})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var jsonData map[string]interface{}
	json.Unmarshal([]byte(string(body)), &jsonData)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "done!",
		"data": jsonData,
	})
}

func ConnectDevice(vins []string, accessToken string, refreshToken string) int {

	base := config.GetTeslaCredential().ProxyUri
	path := "/api/1/vehicles/fleet_telemetry_config"
	url := fmt.Sprintf("%s%s", base, path)

	fieldToStream1 := "Location"
	fieldToStream2 := "BatteryLevel"
	fieldToStream3 := "VehicleSpeed"
	fieldToStream4 := "Odometer"

	telemetryData := TelemetryRequest{
		Config: Config{
			PreferTyped: true,
			Port:        8443,
			Exp:         1770670000,
			AlertTypes:  []string{"service"},
			Fields: map[string]FieldConfig{
				fieldToStream1: {
					ResendIntervalSeconds: 3600,
					MinimumDelta:          1,
					IntervalSeconds:       60,
				},
				fieldToStream2: {
					ResendIntervalSeconds: 3600,
					MinimumDelta:          1,
					IntervalSeconds:       60,
				},
				fieldToStream3: {
					ResendIntervalSeconds: 3600,
					MinimumDelta:          1,
					IntervalSeconds:       60,
				},
				fieldToStream4: {
					ResendIntervalSeconds: 3600,
					MinimumDelta:          1,
					IntervalSeconds:       60,
				},
			},
			CA:       config.GetTeslaCredential().Certificate,
			Hostname: config.GetTeslaCredential().ServerDomain,
		},
		Vins: vins,
	}

	payload, err := json.Marshal(telemetryData)
	if err != nil {
		fmt.Println("Failed to marshal JSON data")
		return 0
	}

	certPEM := config.GetTeslaCredential().Certificate

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(certPEM)); !ok {
		log.Fatal("Failed to append certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		ServerName: config.GetTeslaCredential().ServerDomain,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal("Error creating request")
		return 0
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making request:")
		return 0
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var jsonData Root
	json.Unmarshal([]byte(string(body)), &jsonData)

	fmt.Println("debug1=>", string(body))
	fmt.Println("debug2=>", jsonData)
	fmt.Println("debug3=>", jsonData.Response.UpdatedVehicles)
	fmt.Println("debug5=>", vins[0])

	if jsonData.Response.UpdatedVehicles == 1 {
		return 2
	} else if jsonData.Response.UpdatedVehicles == 0 {
		return 1
	} else {
		fmt.Println("debug4=>", refreshToken)
		_, err := utils.RefreshAuthToken(refreshToken, vins[0])
		if err != nil {
			fmt.Println(err)
		}

		return 0
	}
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func GetDeviceConfigStatus(c *gin.Context) {
	var requestParams RequestConnectParams
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	base := config.GetTeslaCredential().ProxyUri
	url := fmt.Sprintf("%s/api/1/vehicles/%s/fleet_telemetry_config", base, requestParams.Vins[0])

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error creating request:",
			"error": err,
		})
		return
	}

	certPEM := config.GetTeslaCredential().Certificate

	// Create a certificate pool and add your certificate.
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(certPEM)); !ok {
		log.Fatal("Failed to append certificate")
	}

	// Create a custom TLS configuration that uses the certificate pool.
	tlsConfig := &tls.Config{
		RootCAs: certPool,
		// Optionally, if you need to specify the expected server name explicitly:
		ServerName: config.GetTeslaCredential().ServerDomain,
	}

	// Create an HTTP client that uses this TLS configuration.
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", requestParams.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error making request:",
			"error": err,
		})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(string(body)), &jsonData)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "done!",
		"data": jsonData,
	})

}

func GetFleetTelemetryError(c *gin.Context) {
	var requestParams RequestConnectParams
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	base := config.GetTeslaCredential().ProxyUri
	url := fmt.Sprintf("%s/api/1/partner_accounts/fleet_telemetry_errors", base)

	certPEM := config.GetTeslaCredential().Certificate

	// Create a certificate pool and add your certificate.
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(certPEM)); !ok {
		log.Fatal("Failed to append certificate")
	}

	// Create a custom TLS configuration that uses the certificate pool.
	tlsConfig := &tls.Config{
		RootCAs: certPool,
		// Optionally, if you need to specify the expected server name explicitly:
		ServerName: config.GetTeslaCredential().ServerDomain,
	}

	// Create an HTTP client that uses this TLS configuration.
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error creating request:",
			"error": err,
		})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", requestParams.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error making request:",
			"error": err,
		})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(string(body)), &jsonData)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "done!",
		"data": jsonData,
	})
}

func GetFleetStatus(c *gin.Context) {
	var requestParams RequestConnectParams
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	base := config.GetTeslaCredential().ProxyUri
	path := "/api/1/vehicles/fleet_status"
	url := fmt.Sprintf("%s%s", base, path)

	// Build the JSON payload dynamically using the VINs from requestParams
	payload := map[string]interface{}{
		"vins": requestParams.Vins, // Assuming `Vins` is an array of strings in `requestParams`
	}

	jsonStr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	cert, err := tls.X509KeyPair(
		[]byte(config.GetTeslaCredential().ClientCert), // Your client certificate
		[]byte(config.GetTeslaCredential().ClientKey),  // Your client private key
	)

	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}

	caCert := []byte(config.GetTeslaCredential().Certificate) // CA certificate
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to append CA certificate")
	}

	// Create a custom TLS configuration with client certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert}, // Include client certificate
		RootCAs:      caCertPool,              // Include CA certificate
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 30 * time.Second,
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", requestParams.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var jsonData map[string]interface{}
	json.Unmarshal([]byte(string(body)), &jsonData)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "done!",
		"data": jsonData,
	})
}

func GetDeviceLiveData(c *gin.Context) {
	var requestParams RequestConnectParams
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	base := config.GetTeslaCredential().ProxyUri
	url := fmt.Sprintf(base+"/api/1/vehicles/%s/vehicle_data", requestParams.Vins[0])

	certPEM := config.GetTeslaCredential().Certificate

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(certPEM)); !ok {
		log.Fatal("Failed to append certificate")
	}

	// Create a custom TLS configuration that uses the certificate pool.
	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		ServerName: config.GetTeslaCredential().ServerDomain,
	}

	// Create an HTTP client that uses this TLS configuration.
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+requestParams.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error making request:",
			"error": err,
		})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var jsonData map[string]interface{}
	json.Unmarshal([]byte(string(body)), &jsonData)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "done!",
		"data": jsonData,
	})
}

func UpdateDeviceInfo(c *gin.Context) {
	var deviceInfoParams DeviceInfoParams
	if err := c.ShouldBindJSON(&deviceInfoParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload: " + err.Error()})
		return
	}

	var payload Payload
	payload.ShareInfo = make(map[string]bool) // Initialize ShareInfo map
	var deviceList []VehicleInfo

	for _, device := range deviceInfoParams.DeviceList {
		if device.Vin != "" {
			base := config.GetTeslaCredential().ProxyUri
			url := fmt.Sprintf(base+"/api/1/vehicles/%s/vehicle_data", device.Vin)

			certPEM := config.GetTeslaCredential().Certificate
			certPool := x509.NewCertPool()
			if ok := certPool.AppendCertsFromPEM([]byte(certPEM)); !ok {
				log.Fatal("Failed to append certificate")
			}

			// TLS Configuration
			tlsConfig := &tls.Config{
				RootCAs:    certPool,
				ServerName: config.GetTeslaCredential().ServerDomain,
			}

			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: tlsConfig,
				},
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg":   "Error creating request",
					"error": err.Error(),
				})
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+deviceInfoParams.AccessToken)

			resp, err := client.Do(req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg":   "Error making request",
					"error": err.Error(),
				})
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			var vehicleInfoParams VehicleInfoParams
			if err := json.Unmarshal(body, &vehicleInfoParams); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg":   "Error parsing vehicle data",
					"error": err.Error(),
				})
				return
			}

			// Populate VehicleInfo
			vehicleInfo := VehicleInfo{
				Vin:         device.Vin,
				CarType:     vehicleInfoParams.Response.VehicleConfig.CarType,
				DeviceName:  device.DeviceName,
				VehicleID:   strconv.Itoa(vehicleInfoParams.Response.VehicleID),
				ShareAbi:    deviceInfoParams.ShareStatus["abi_insurance"],
				ShareTintAi: deviceInfoParams.ShareStatus["tint_ai"],
				Color:       vehicleInfoParams.Response.Color,
			}

			deviceList = append(deviceList, vehicleInfo)
		}
	}

	// Populate Payload
	payload.AccessToken = deviceInfoParams.AccessToken
	payload.RefreshToken = deviceInfoParams.RefreshToken
	payload.Email = deviceInfoParams.Email
	payload.DeviceList = deviceList
	payload.ShareInfo["abi_insurance"] = deviceInfoParams.ShareStatus["abi_insurance"]
	payload.ShareInfo["tint_ai"] = deviceInfoParams.ShareStatus["tint_ai"]

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	url := config.GetTeslaCredential().TestServerUri + "/api/tesla_device_signup"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Print the response
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	c.JSON(http.StatusOK, gin.H{
		"data": payload,
		"msg":  "done!",
	})
}
