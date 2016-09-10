package account

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jchorl/financejc/constants"

	"github.com/Sirupsen/logrus"
)

var InvalidCurrency = errors.New("Currency is invalid")

type Account struct {
	Id       int    `json:"id,omitempty" description:"Id of the account"`
	Name     string `json:"name" description:"Name of the account"`
	Currency string `json:"currency" description:"Currency for the account"`
	User     int    `json:"-"`
}

func Get(c context.Context, userId string) ([]*Account, error) {
	accounts := make([]*Account, 0)
	db := c.Value(constants.CTX_DB).(sql.DB)
	rows, err := db.Query("SELECT * FROM accounts WHERE user = $1", userId)
	if err != nil {
		logrus.WithField("Error", err).Error("failed to fetch accounts")
		return nil, err
	}

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account); err != nil {
			logrus.WithField("Error", err).Error("failed to scan into account")
			return nil, err
		}

		accounts = append(accounts, &account)
	}
	if err := rows.Err(); err != nil {
		logrus.WithField("Error", err).Error("failed to get accounts from rows")
		return nil, err
	}

	return accounts, nil
}

func New(c context.Context, userId string, account *Account) (*Account, error) {
	_, valid := constants.CurrencyCodeToName[account.Currency]
	if !valid {
		return nil, InvalidCurrency
	}

	db := c.Value(constants.CTX_DB).(sql.DB)
	var id int
	err := db.QueryRow("INSERT INTO accounts(name, currency) VALUES($1, $2) RETURNING id", account.Name, account.Currency).Scan(&id)
	if err != nil {
		logrus.WithField("Error", err).Errorf("failed to insert account row")
		return nil, err
	}

	account.Id = id
	return account, nil
}

func Update(c context.Context, account *Account) (*Account, error) {
	_, valid := constants.CurrencyCodeToName[account.Currency]
	if !valid {
		return nil, InvalidCurrency
	}

	db := c.Value(constants.CTX_DB).(sql.DB)
	_, err := db.Exec("UPDATE accounts SET name = $1, currency = $2 WHERE id = $3", account.Name, account.Currency, account.Id)
	if err != nil {
		logrus.WithField("Error", err).Errorf("failed to update account row")
		return nil, err
	}

	return account, nil
}

func Delete(c context.Context, accountId string) error {
	db := c.Value(constants.CTX_DB).(sql.DB)

	_, err := db.Exec("DELETE FROM accounts WHERE id = $1", accountId)
	if err != nil {
		logrus.WithField("Error", err).Errorf("could not delete account")
		return err
	}

	return nil
}
