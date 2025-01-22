package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	config "tesla_server/condig"

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

	base := "https://fleet-api.prd.na.vn.cloud.tesla.com"
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
				"<field_to_stream>": map[string]interface{}{
					"resend_interval_seconds": 3600,
					"minimum_delta":           1,
					"interval_seconds":        1800,
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
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TESLA_AUTH_TOKEN")))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)

	fmt.Println("resp=>", resp)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println(string(body))
}
