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

	"github.com/gin-gonic/gin"
)

type RequestConnectParams struct {
	Vins         []string `json:"vins"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
}

type Config struct {
	Port   int    `json:"port"`
	Exp    int    `json:"exp"`
	Fields Fields `json:"fields"`
	CA     string `json:"ca"`
	Host   string `json:"hostname"`
}

type Fields struct {
	Location    Interval `json:"Location"`
	ChargeState Interval `json:"ChargeState"`
}

type Interval struct {
	IntervalSeconds int `json:"interval_seconds"`
}

type TelemetryData struct {
	Config Config   `json:"config"`
	VINs   []string `json:"vins"`
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

	// Construct headers
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", requestParams.AccessToken),
	}

	// Construct data
	data := TelemetryData{
		Config: Config{
			Port: 8443,
			Exp:  1750000000,
			Fields: Fields{
				Location:    Interval{IntervalSeconds: 500},
				ChargeState: Interval{IntervalSeconds: 500},
			},
			CA:   config.GetTeslaCredential().Certificate,
			Host: "t3slaapi.moovetrax.com",
		},
		VINs: requestParams.Vins,
	}

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal JSON data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Failed to marshal JSON data:",
			"error": err,
		})
		return
	}

	// Set up TLS configuration
	tlsConfig := &tls.Config{
		RootCAs: nil, // Optionally load CA from caCert if needed
	}

	caCert := config.GetTeslaCredential().TlsCertificate
	if caCert != "" {
		// Create a new certificate pool
		roots := x509.NewCertPool()

		// Append the CA certificate to the certificate pool
		if ok := roots.AppendCertsFromPEM([]byte(caCert)); !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":   "Failed to parse CA certificate:",
				"error": err,
			})
		}

		// Assign the certificate pool to the TLS config
		tlsConfig.RootCAs = roots
	}
	// Create HTTP client with custom TLS configuration
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Create POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Failed to create request:",
			"error": err,
		})
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Failed to send request:",
			"error": err,
		})
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Failed to read response body:",
			"error": err,
		})
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Response:", string(body))
		c.JSON(http.StatusOK, gin.H{
			"msg":   "successfully:",
			"error": string(body),
		})

	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Failed to read response body:",
			"error": string(body),
		})
	}
}
