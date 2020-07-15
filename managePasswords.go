package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// BulkPasswordReset sets all pws
func BulkPasswordReset(users []UserDeprecated) error {
	for _, u := range users {

		resetPwURI := viper.GetString("resetPwUri")

		fmt.Println("reset for ", u.Email, " to ", resetPwURI)
		u := User{
			Email: u.Email,
		}
		reqBody, err := json.Marshal(u)
		if err != nil {
			print(err)
			return err
		}
		req, err := http.NewRequest("POST",
			resetPwURI,
			bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(body))

		// pause to give deprecated api some rest
		duration := time.Second
		time.Sleep(duration)

	}
	return nil
}
