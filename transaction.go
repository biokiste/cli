package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func getUserTransactions(db *sql.DB, dbOld *sql.DB, user UserDeprecated) []TransactionDeprecated {
	results, err := dbOld.Query(`
	SELECT transactions.id, amount, transactions.created_at, firstname, lastname, transactions.status, transactions.reason, category_id, transactions_category.type
	FROM transactions
	LEFT JOIN transactions_category ON transactions.category_id = transactions_category.id
	LEFT JOIN users ON transactions.user_id = users.id
	WHERE users.id = ?
	ORDER BY transactions.created_at desc
	`, user.ID)
	if err != nil {
		panic(err)
	}

	defer results.Close()
	var transactions []TransactionDeprecated
	for results.Next() {
		var transaction TransactionDeprecated

		err = results.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.CreatedAt,
			&transaction.FirstName,
			&transaction.LastName,
			&transaction.Status,
			&transaction.Reason,
			&transaction.CategoryID,
			&transaction.Type)
		if err != nil {
			panic(err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}

func getUserIDMapping(db *sql.DB, user UserDeprecated) int {
	var u User

	query := fmt.Sprintf(`	
	SELECT 	ID, AdditionalInfos 
	FROM Users	
	WHERE AdditionalInfos = %d`, user.ID)

	err := db.QueryRow(query).Scan(&u.ID, &u.AdditionalInfos)
	if err != nil {
		return 0
	}

	// do not migrate transactions of former users
	if user.State == 4 {
		return 0
	}

	return u.ID
}

// AddUserTransaction sends user transaction
func AddUserTransaction(db *sql.DB, dbOld *sql.DB, users []UserDeprecated) error {

	token := viper.GetString("token")
	apiBaseURL := viper.GetString("api_base_url")
	transactionTypes := viper.GetStringSlice("transaction_types")
	transactionStates := viper.GetStringSlice("transaction_states")

	for _, usr := range users {
		transactions := getUserTransactions(db, dbOld, usr)
		newUserID := getUserIDMapping(db, usr)

		if newUserID != 0 {
			for i, t := range transactions {

				if t.Status > 0 {
					newTransaction := Transaction{
						Amount:        t.Amount,
						Type:          transactionTypes[t.CategoryID-1],
						State:         transactionStates[t.Status-1],
						UserID:        newUserID,
						CreatedAt:     t.CreatedAt,
						CreatedBy:     newUserID,
						UpdateComment: t.Reason,
					}
					fmt.Fprintln(os.Stdout, "user", usr.Email, "transaction ", i, "of", len(transactions))

					if t.Amount != 0 {
						reqBody, err := json.Marshal(newTransaction)
						if err != nil {
							print(err)
						}
						req, err := http.NewRequest("POST", apiBaseURL+"api/transactions", bytes.NewBuffer(reqBody))
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

			}
		}

	}
	return nil
}
