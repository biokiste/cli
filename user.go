package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func prepareUser(users []UserDeprecated) []User {
	var newUsers []User
	for _, usr := range users {
		userMobile := strings.ReplaceAll(usr.Mobile, " ", "")

		var state string
		switch usr.State {
		case 2:
			state = "active"
		case 4:
			state = "former"
		case 6:
			state = "paused"
		default:
			state = "nil"
		}

		var birthday string
		if usr.DateOfBirth == "0000-00-00" {
			birthday = "2222-12-12" // fallback for invalid date
		} else {
			birthday = usr.DateOfBirth
		}

		var leavingDate string
		if usr.State == 4 {
			leavingDate = usr.DateOfExit
		} else {
			leavingDate = ""
		}

		newUser := User{
			FirstName:       usr.Firstname,
			LastName:        usr.Lastname,
			Email:           usr.Email,
			Phone:           userMobile,
			Street:          usr.Street,
			StreetNumber:    "0815", // TODO: split Street to get Number ?
			Zip:             usr.ZIP,
			Country:         "Germany",
			Birthday:        birthday,
			EntranceDate:    usr.DateOfEntry,
			LeavingDate:     leavingDate,
			AdditionalInfos: strconv.Itoa(usr.ID), // save oldID for further transaction inserts
			LastActivityAt:  usr.LastLogin,
			CreatedAt:       usr.CreatedAt,
			UpdatedAt:       usr.UpdatedAt,
			CreatedBy:       174,
			Password:        "2020_b!ok!ste_" + userMobile, // initial password
			State:           state,
		}

		if state != "nil" {
			newUsers = append(newUsers, newUser)
		}

	}
	return newUsers

}

// AddUserReq sends user creation
func AddUserReq(users []UserDeprecated) error {
	token := viper.GetString("token")
	apiBaseURL := viper.GetString("api_base_url")
	newUsers := prepareUser(users)

	for _, addUser := range newUsers {
		fmt.Println(addUser)
		reqBody, err := json.Marshal(addUser)
		if err != nil {
			print(err)
		}
		req, err := http.NewRequest("POST", apiBaseURL+"api/users", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			print(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			print(err)
		}

		if resp.StatusCode != 200 {
			fmt.Println(resp.StatusCode, string(body))
			return errors.New(string(resp.StatusCode))
		}

		fmt.Println("resp ", string(body))
	}

	return nil
}
