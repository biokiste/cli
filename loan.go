package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

func getUsers(db *sql.DB) []User {
	results, err := db.Query(`
	SELECT
		ID,			
		COALESCE(AdditionalInfos, '') as AdditionalInfos,
		State
	FROM Users`)

	if err != nil {
		panic(err)
	}

	var users []User
	for results.Next() {
		var user User

		err = results.Scan(
			&user.ID,
			&user.AdditionalInfos,
			&user.State,
		)
		if err != nil {
			panic(err)
		}
		users = append(users, user)

		defer results.Close()

	}
	return users
}

func getUserLoan(dbOld *sql.DB, user User) UserDeprecated {
	row := dbOld.QueryRow(`
	SELECT  credit, credit_date, credit_comment
	FROM users	
	WHERE id = ?	
	`, user.AdditionalInfos)

	var usr UserDeprecated

	row.Scan(
		&usr.Credit,
		&usr.CreditDate,
		&usr.CreditComment,
	)

	return usr
}

// AddUserLoan sends user loan
func AddUserLoan(db *sql.DB, dbOld *sql.DB) error {

	token := viper.GetString("token")
	apiBaseURL := viper.GetString("api_base_url")
	users := getUsers(db) // get users from migrated db

	for _, usr := range users {
		depUsr := getUserLoan(dbOld, usr)
		loan := Loan{
			UserID:        usr.ID,
			State:         usr.State,
			Amount:        depUsr.Credit,
			CreatedAt:     depUsr.CreditDate,
			UpdateComment: depUsr.CreditComment,
			CreatedBy:     usr.ID,
		}

		if depUsr.Credit > 0 {
			reqBody, err := json.Marshal(loan)
			if err != nil {
				print(err)
			}
			req, err := http.NewRequest("POST", apiBaseURL+"api/loans", bytes.NewBuffer(reqBody))
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
				return errors.New(string(resp.StatusCode))
			}

			fmt.Println("resp ", string(body))
		}

	}
	return nil
}
