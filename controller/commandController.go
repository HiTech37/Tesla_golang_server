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
	"strings"
	config "tesla_server/config"

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

	jsonData := map[string]interface{}{}
	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	caCert := []byte(config.GetTeslaCredential().Certificate) // Replace with your CA certificate
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to append CA certificate")
	}

	// Create a custom TLS configuration
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", requestParams.AccessToken))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resData := strings.TrimSpace(string(body))
	if strings.Contains(resData, `"error":"token expired (401)"`) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "token updated",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resData,
	})
}
