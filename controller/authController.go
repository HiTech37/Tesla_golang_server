package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	config "tesla_server/config"

	"github.com/gin-gonic/gin"
)

type RequestData struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	Audience     string `json:"audience"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

type Response struct {
	AccessToken string `json:"access_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GetRoutes handles GET requests to "/api/get-routes"
func GetTeslaSigninURI(c *gin.Context) {

	credentail := config.GetTeslaCredential()
	signinUri := "https://auth.tesla.com/oauth2/v3/authorize?response_type=code&locale=en-US&prompt=login&client_id=" + credentail.ClientID + "&redirect_uri=" + url.QueryEscape(credentail.CallbackUri) + "&scope=" + credentail.DataScope + "&state=db4af3f87"

	c.Writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	encoder.SetEscapeHTML(false)

	// Create the response map
	response := gin.H{
		"signin_uri": signinUri,
	}

	// Encode the response with the custom encoder
	if err := encoder.Encode(response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to encode response",
		})
		return
	}
}

func RequestAuth(c *gin.Context) {

	code := c.Query("code")
	externalURL := "https://auth.tesla.com/oauth2/v3/token"

	credential := config.GetTeslaCredential()

	data := RequestData{
		GrantType:    "client_credentials",
		ClientID:     credential.ClientID,
		ClientSecret: credential.SecretKey,
		Audience:     "https://fleet-api.prd.na.vn.cloud.tesla.com",
		Scope:        credential.DataScope,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	req, err := http.NewRequest("POST", externalURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// If successful, print the response status and body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var responseObj Response
	err = json.Unmarshal(body, &responseObj)
	if err != nil {
		fmt.Println("Error unmarshaling response body:", err)
		return
	}

	res := RegisterPublicKey(responseObj.AccessToken)

	if res {

		// URL for the Tesla OAuth token endpoint
		tokenURL := "https://auth.tesla.com/oauth2/v3/token"

		// Prepare form data
		formData := url.Values{}
		formData.Set("grant_type", "authorization_code")
		formData.Set("client_id", credential.ClientID)
		formData.Set("client_secret", credential.SecretKey)
		formData.Set("code", code)
		formData.Set("audience", "https://fleet-api.prd.eu.vn.cloud.tesla.com")
		formData.Set("redirect_uri", credential.CallbackUri)

		// Create the POST request
		req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(formData.Encode()))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		// Set headers
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
			return
		}
		defer resp.Body.Close()

		// Read and return the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		var tokenResponse TokenResponse
		err = json.Unmarshal(body, &tokenResponse)
		if err != nil {
			fmt.Println("Error unmarshaling response body:", err)
			return
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(c.Writer)
		encoder.SetEscapeHTML(false)

		response := gin.H{
			"access_token":  tokenResponse.AccessToken,
			"refresh_token": tokenResponse.RefreshToken,
		}
		// Encode the response with the custom encoder
		if err := encoder.Encode(response); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to encode response",
			})
			return
		}

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{})
	}

}

func RegisterPublicKey(accessToken string) bool {
	credential := config.GetTeslaCredential()

	url := "https://fleet-api.prd.na.vn.cloud.tesla.com/api/1/partner_accounts"

	// Create the request body
	data := map[string]string{
		"domain": credential.ServerDomain,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

func GetAllVehicles(c *gin.Context) {
	accessToken := c.Query("access_token")
	tesla_api_uri := "https://fleet-api.prd.na.vn.cloud.tesla.com"
	url := tesla_api_uri + "/api/1/vehicles"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	encoder := json.NewEncoder(c.Writer)
	encoder.SetEscapeHTML(false)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}
	// Return the response body to the client
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

func GetAccessToken(refreshToken string) {

}
