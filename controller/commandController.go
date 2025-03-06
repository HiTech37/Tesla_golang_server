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
	"strings"
	config "tesla_server/config"
	"tesla_server/utils"

	"github.com/gin-gonic/gin"
)

type RequestCommandParams struct {
	Command      string `json:"command"`
	Vin          string `json:"vin"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func HandleCommand(c *gin.Context) {
	var requestParams RequestCommandParams
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var url string

	switch requestParams.Command {
	case "Unlock":
		url = config.GetTeslaCredential().ProxyUri + fmt.Sprintf("/api/1/vehicles/%s/command/door_unlock", requestParams.Vin)
	case "Lock":
		url = config.GetTeslaCredential().ProxyUri + fmt.Sprintf("/api/1/vehicles/%s/command/door_lock", requestParams.Vin)
	case "Light":
		url = config.GetTeslaCredential().ProxyUri + fmt.Sprintf("/api/1/vehicles/%s/command/flash_lights", requestParams.Vin)
	case "HonkHorn":
		url = config.GetTeslaCredential().ProxyUri + fmt.Sprintf("/api/1/vehicles/%s/command/honk_horn", requestParams.Vin)
	default:
		url = ""
	}

	if url == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unsupported command!",
		})
		return
	}

	resData, err := SendCommand(url, requestParams.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resData = strings.TrimSpace(resData)
	if strings.Contains(resData, `"error":"token expired (401)"`) {
		var teslaAuthToken utils.TeslaAuthToken
		teslaAuthToken, err := utils.RefreshAuthToken(requestParams.RefreshToken, requestParams.Vin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		req, err := SendCommand(url, teslaAuthToken.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": req,
			"msg":  "done",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resData,
		"msg":  "done",
	})
}

func SendCommand(url string, accessToken string) (string, error) {

	jsonData := map[string]interface{}{}
	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}

	caCert := []byte(config.GetTeslaCredential().Certificate) // Replace with your CA certificate
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to append CA certificate")
	}

	// Create a custom TLS configuration
	tlsConfig := &tls.Config{
		RootCAs:    caCertPool,
		ServerName: config.GetTeslaCredential().ServerDomain,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
