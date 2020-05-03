package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

func intToBool(i int) bool {
	if i == 1 {
		return true
	}
	return false
}

// CreateGroups creates groups
func CreateGroups(dbOld *sql.DB) error {
	token := viper.GetString("token")
	apiBaseURL := viper.GetString("api_base_url")

	query := "SELECT name FROM groups"
	results, err := dbOld.Query(query)

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer results.Close()

	var groupsDeprecated []GroupDeprecated
	for results.Next() {
		var g GroupDeprecated
		err := results.Scan(
			&g.Name,
		)
		if err != nil {
			fmt.Println(err)
			return err
		}
		groupsDeprecated = append(groupsDeprecated, g)
	}

	for _, g := range groupsDeprecated {
		group := Group{
			GroupKey:  g.Name,
			Email:     "mailingliste@biokiste.org",
			CreatedBy: -1,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		reqBody, err := json.Marshal(group)
		if err != nil {
			print(err)
			return err
		}
		req, err := http.NewRequest("POST",
			apiBaseURL+"api/groups",
			bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
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

	}

	return nil
}

// AddUserToGroups sets user group relations
func AddUserToGroups(dbOld *sql.DB, db *sql.DB) error {
	results, err := db.Query(`SELECT  ID, AdditionalInfos	FROM users`)
	if err != nil {
		fmt.Println(err)
	}

	defer results.Close()

	var users []User
	for results.Next() {
		var usr User

		err := results.Scan(
			&usr.ID,
			&usr.AdditionalInfos,
		)
		if err != nil {
			fmt.Println(err)

		}
		users = append(users, usr)
	}

	for _, u := range users {
		results, err := dbOld.Query(`
		SELECT  group_id, position_id	
		FROM groups_users
		WHERE active = 1
		AND user_id = ?
		`, u.AdditionalInfos)

		if err != nil {
			fmt.Println(err)
		}
		defer results.Close()

		for results.Next() {
			var uDep UserGroupDeprecated
			err := results.Scan(
				&uDep.GroupID,
				&uDep.PositionID,
			)
			if err != nil {
				fmt.Println(err)

			}

			err = addUserToGroup(u.ID, uDep)

		}

	}

	return nil
}

// AddUserToGroup sends user loan
func addUserToGroup(usrID int, uDep UserGroupDeprecated) error {

	token := viper.GetString("token")
	apiBaseURL := viper.GetString("api_base_url")

	userGroup := UserGroup{
		GroupID:  uDep.GroupID,
		IsLeader: intToBool(uDep.PositionID),
	}

	reqBody, err := json.Marshal(userGroup)
	if err != nil {
		fmt.Println(err)
	}

	userID := strconv.Itoa(usrID)

	req, err := http.NewRequest("POST",
		apiBaseURL+"api/users/"+userID+"/groups",
		bytes.NewBuffer(reqBody))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
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

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode, string(body))
		return errors.New(string(resp.StatusCode))
	}

	fmt.Println("resp ", string(body))

	return nil
}
