package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

type userCreationResponse struct {
	ID string `json:"id"`
}

// AddUserReq sends user creation
func AddUserReq(user User) (string, error) {
	token := viper.GetString("token")
	apiBaseURL := viper.GetString("api_base_url")

	reqBody, err := json.Marshal(user)
	if err != nil {
		print(err)
	}
	req, err := http.NewRequest("POST", apiBaseURL+"api/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode, string(body))
		return string(body), errors.New(string(resp.StatusCode))
	}

	var u userCreationResponse
	json.NewDecoder(resp.Body).Decode(&u)

	// TODO: id should be assign
	fmt.Println("resp ", string(body))
	fmt.Println("user ", u)

	return u.ID, nil
}
