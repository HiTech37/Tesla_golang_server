package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
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
	url := "https://auth.tesla.com/oauth2/v3/token"
	method := "POST"
	clientID := config.GetTeslaCredential().ClientID
	scope := config.GetTeslaCredential().DataScope
	payload := strings.NewReader(fmt.Sprintf("grant_type=refresh_token&client_id=%s&scope=%s&refresh_token=%s", clientID, scope, refreshToken))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return teslaAuthToken, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return teslaAuthToken, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return teslaAuthToken, err
	}

	err = json.Unmarshal(body, &teslaAuthToken)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return teslaAuthToken, err
	}

	if teslaAuthToken.AccessToken != "" && teslaAuthToken.RefreshToken != "" && vin != "" {
		err := model.UpdateDeviceAuthTokensbyVin(teslaAuthToken.AccessToken, teslaAuthToken.RefreshToken, vin)
		if err != nil {
			return teslaAuthToken, err
		}
		return teslaAuthToken, nil
	}

	return teslaAuthToken, nil
}
