package controller

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	config "tesla_server/config"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestConnectParams struct {
	Vins         []string `json:"vins"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
}

type TelemetryRequest struct {
	Config struct {
		Port   int `json:"port"`
		Exp    int `json:"exp"`
		Fields struct {
			Location struct {
				IntervalSeconds int `json:"interval_seconds"`
			} `json:"Location"`
			ChargeState struct {
				IntervalSeconds int `json:"interval_seconds"`
			} `json:"ChargeState"`
		} `json:"fields"`
		CA       string `json:"ca"`
		Hostname string `json:"hostname"`
	} `json:"config"`
	VINs []string `json:"vins"`
}

func ConnectDevice(c *gin.Context) {
	var requestParams RequestConnectParams
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	base := config.GetTeslaCredential().ProxyUri
	path := "/api/1/vehicles/fleet_telemetry_config"
	url := fmt.Sprintf("%s%s", base, path)

	telemetryData := TelemetryRequest{}
	telemetryData.Config.Port = 8443
	telemetryData.Config.Exp = 1750000000
	telemetryData.Config.Fields.Location.IntervalSeconds = 500
	telemetryData.Config.Fields.ChargeState.IntervalSeconds = 500
	telemetryData.Config.CA = config.GetTeslaCredential().Certificate

	telemetryData.Config.Hostname = config.GetTeslaCredential().ServerDomain
	telemetryData.VINs = requestParams.Vins

	payload, err := json.Marshal(telemetryData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Failed to marshal JSON data:",
			"error": err,
		})
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

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "failed to create HTTP request",
			"error": err,
		})
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", requestParams.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":      "failed to send HTTP request:",
			"error":    err,
			"jsonData": telemetryData,
			"payload":  payload,
		})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonData map[string]interface{}
	json.Unmarshal([]byte(string(body)), &jsonData)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "done!",
		"data": jsonData,
	})
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", requestParams.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error making request:",
			"error": err,
		})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error making request:",
			"error": err,
		})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
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

	body, err := ioutil.ReadAll(resp.Body)
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
