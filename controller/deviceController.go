package controller

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	telemetryData.Config.Fields.Location.IntervalSeconds = 5
	telemetryData.Config.Fields.ChargeState.IntervalSeconds = 5
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

	// Create a custom TLS configuration with CA cert
	rootCAs := x509.NewCertPool()
	if ok := rootCAs.AppendCertsFromPEM([]byte(config.GetTeslaCredential().TlsCertificate)); !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "failed to append CA certificate",
			"error": err,
		})
		return
	}
	tlsConfig := &tls.Config{
		RootCAs: rootCAs,
	}

	// Create HTTPS client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 10 * time.Second,
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

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "failed to send HTTP request:",
			"error": err,
		})
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	c.JSON(http.StatusOK, gin.H{
		"msg":  "failed to send HTTP request:",
		"data": string(body),
	})
}
