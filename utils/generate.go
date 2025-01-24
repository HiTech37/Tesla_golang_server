package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"tesla_server/config"
	"tesla_server/model"
	"time"
)

func GenerateRandomString(strLength int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	// Create a slice of random bytes
	b := make([]byte, strLength)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

type TeslaAuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func RefreshAuthToken(refreshToken string, vin string) (TeslaAuthToken, error) {
	var teslaAuthToken TeslaAuthToken
	// URL for the request
	apiURL := "https://fleet-auth.prd.vn.cloud.tesla.com/oauth2/v3/token"

	// Form data
	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("client_id", config.GetTeslaCredential().ClientID)
	formData.Set("refresh_token", refreshToken)

	// Create a new POST request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return teslaAuthToken, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return teslaAuthToken, err
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return teslaAuthToken, err
	}

	err = json.Unmarshal(body, &teslaAuthToken)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return teslaAuthToken, err
	}

	fmt.Println("debug1=>", teslaAuthToken)
	if teslaAuthToken.AccessToken != "" && teslaAuthToken.RefreshToken != "" && vin != "" {
		err := model.UpdateDeviceAuthTokensbyVin(teslaAuthToken.AccessToken, teslaAuthToken.RefreshToken, vin)
		if err != nil {
			return teslaAuthToken, err
		}
		return teslaAuthToken, nil
	}

	return teslaAuthToken, nil
}
