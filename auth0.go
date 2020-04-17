package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

const users = "api/v2/users/"

// Auth0ErrorResponse represents error response
type Auth0ErrorResponse struct {
	StatusCode int
	Error      string `json:"error"`
	Message    string `json:"message"`
}

// Auth0Bearer represents token object
type Auth0Bearer struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Auth0User represents Auth0 User data
type Auth0User struct {
	Connection  string `json:"connection"`
	UserID      string `json:"user_id,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	LastLogin   string `json:"last_login,omitempty"`
	VerifyEmail bool   `json:"verify_email"`
}

func getToken() (string, error) {
	auth0URI := viper.GetString("auth0URI")
	clientID := viper.GetString("clientId")
	clientSecret := viper.GetString("clientSecret")
	clientAudience := viper.GetString("audience")

	if len(auth0URI) == 0 || len(clientID) == 0 || len(clientSecret) == 0 || len(clientAudience) == 0 {
		err := errors.New("no credentials for auth provider found")
		return "", err
	}
	params := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s&audience=%s", clientID, clientSecret, clientAudience)
	payload := strings.NewReader(params)
	url := fmt.Sprintf("%soauth/token", auth0URI)
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		var errRes Auth0ErrorResponse
		json.NewDecoder(res.Body).Decode(&errRes)
		errRes.StatusCode = res.StatusCode
		return "", errors.New(errRes.Message)
	}

	var bearer Auth0Bearer
	json.NewDecoder(res.Body).Decode(&bearer)

	return bearer.AccessToken, nil
}

// DeleteAuth0User deletes user at auth0
func DeleteAuth0User(userID string) error {
	auth0URI := viper.GetString("auth0URI")
	apikey, err := getToken()

	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", auth0URI+users+userID, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apikey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		return err
	}
	fmt.Println("User " + userID + " removed")
	return nil
}
