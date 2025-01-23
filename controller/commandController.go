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

	var resData string
	var resErr error

	switch requestParams.Command {
	case "Unlock":
		resData, resErr = SendCommandUnlock(requestParams.AccessToken, requestParams.Vin)
	case "Lock":
		resData, resErr = SendCommandLock(requestParams.AccessToken, requestParams.Vin)
	case "Light":
		resData, resErr = SendCommandLight(requestParams.AccessToken, requestParams.Vin)
	case "HonkHorn":
		resData, resErr = SendCommandHonkHorn(requestParams.AccessToken, requestParams.Vin)
	default:
		resErr = fmt.Errorf("invalid command: %s", requestParams.Command)
	}

	fmt.Println(requestParams.Command)
	fmt.Println(requestParams.Vin)
	fmt.Println(resData)
	fmt.Println(resErr)
	if resErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": resErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  requestParams.AccessToken,
		"refreshToken": requestParams.RefreshToken,
		"data":         resData,
	})
}

func SendCommandUnlock(accessToken string, vehicleTag string) (string, error) {
	// base := "https://fleet-api.prd.na.vn.cloud.tesla.com"
	base := config.GetTeslaCredential().ProxyUri
	path := fmt.Sprintf("/api/1/vehicles/%s/command/door_unlock", vehicleTag)
	url := base + path

	// Prepare JSON payload
	jsonData := map[string]interface{}{}
	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return "", err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	return string(body), nil
}

func SendCommandLock(accessToken string, vehicleTag string) (string, error) {
	base := "https://fleet-api.prd.na.vn.cloud.tesla.com"
	path := fmt.Sprintf("/api/1/vehicles/%s/command/door_lock", vehicleTag) // Adjusting the path for locking the door
	url := base + path

	// Prepare JSON payload
	jsonData := map[string]interface{}{}
	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return "", err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read and return the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	return string(body), nil
}

func SendCommandLight(accessToken string, vehicleTag string) (string, error) {
	base := "https://fleet-api.prd.na.vn.cloud.tesla.com"
	path := fmt.Sprintf("/api/1/vehicles/%s/command/flash_lights", vehicleTag) // Corrected to use vehicleTag in the path
	url := base + path

	// Prepare JSON payload
	jsonStr := `{}`

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken)) // Using the accessToken parameter

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read and return the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	return string(body), nil
}

func SendCommandHonkHorn(accessToken string, vehicleTag string) (string, error) {
	base := "https://fleet-api.prd.na.vn.cloud.tesla.com"
	path := fmt.Sprintf("/api/1/vehicles/%s/command/honk_horn", vehicleTag) // Use vehicleTag in the path
	url := base + path

	// Prepare JSON payload
	jsonStr := `{}`

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken)) // Use the accessToken parameter

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read and return the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}
	return string(body), nil
}
