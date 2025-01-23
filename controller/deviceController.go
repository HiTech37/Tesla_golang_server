package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	config "tesla_server/config"

	"github.com/gin-gonic/gin"
)

type RequestConnectParams struct {
	Vins         []string `json:"vins"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
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

	// Define the data structure
	data := map[string]interface{}{
		"config": map[string]interface{}{
			"prefer_typed": true,
			"port":         config.GetTeslaCredential().Port,
			"exp":          1704067200,
			"alert_types":  []string{"service"},
			"fields": map[string]interface{}{
				"Location": map[string]interface{}{
					"resend_interval_seconds": 10,
					"minimum_delta":           1,
					"interval_seconds":        5,
				},
			},
			"ca":       config.GetTeslaCredential().Certificate,
			"hostname": config.GetTeslaCredential().ServerDomain,
		},
		"vins": requestParams.Vins,
	}

	// Marshal the data structure into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error marshaling JSON:",
			"error": err,
		})
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error creating request:",
			"error": err,
		})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", requestParams.AccessToken))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making request:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error making request:",
			"error": err,
			"ca":    config.GetTeslaCredential().Certificate,
		})
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error reading response body:",
			"error": err,
		})
		return
	}
	encoder := json.NewEncoder(c.Writer)
	encoder.SetEscapeHTML(false)

	response := gin.H{
		"status": resp.StatusCode,
		"data":   string(body),
	}

	if err := encoder.Encode(response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to encode response",
		})
		return
	}
}
