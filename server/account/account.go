package account

import (
	"context"
	"errors"

	"github.com/jchorl/financejc/constants"
	"github.com/jchorl/financejc/server/util"

	"github.com/Sirupsen/logrus"
)

var InvalidCurrency = errors.New("Currency is invalid")

type Account struct {
	Id          int     `json:"id,omitempty" description:"Id of the account"`
	Name        string  `json:"name" description:"Name of the account"`
	Currency    string  `json:"currency" description:"Currency for the account"`
	User        int     `json:"-"`
	FutureValue float64 `json:"futureValue" description:"Future value of the account"`
}

func Get(c context.Context, userId int) ([]*Account, error) {
	accounts := make([]*Account, 0)
	db := c.Value(constants.CTX_DB).(util.DB)
	rows, err := db.Query("SELECT a.id, a.name, a.currency, a.userId, SUM(t.amount) FROM accounts a JOIN transactions t on t.account=a.id WHERE a.userId = $1 GROUP BY a.id", userId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":   err,
			"User ID": userId,
		}).Error("failed to fetch accounts")
		return nil, err
	}

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.Id, &account.Name, &account.Currency, &account.User, &account.FutureValue); err != nil {
			logrus.WithFields(logrus.Fields{
				"Error":   err,
				"User ID": userId,
			}).Error("failed to scan into account")
			return nil, err
		}

		accounts = append(accounts, &account)
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":   err,
			"User ID": userId,
		}).Error("failed to get accounts from rows")
		return nil, err
	}

	return accounts, nil
}

func New(c context.Context, account *Account) (*Account, error) {
	_, valid := constants.CurrencyCodeToName[account.Currency]
	if !valid {
		return nil, InvalidCurrency
	}

	db := c.Value(constants.CTX_DB).(util.DB)
	var id int
	err := db.QueryRow("INSERT INTO accounts(name, currency, userId) VALUES($1, $2, $3) RETURNING id", account.Name, account.Currency, account.User).Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":   err,
			"Account": account,
		}).Errorf("failed to insert account row")
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

	db := c.Value(constants.CTX_DB).(util.DB)
	_, err := db.Exec("UPDATE accounts SET name = $1, currency = $2 WHERE id = $3", account.Name, account.Currency, account.Id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":   err,
			"Account": account,
		}).Errorf("failed to update account row")
		return nil, err
	}

	return account, nil
}

func Delete(c context.Context, accountId string) error {
	db := c.Value(constants.CTX_DB).(util.DB)

	_, err := db.Exec("DELETE FROM accounts WHERE id = $1", accountId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":      err,
			"Account ID": accountId,
		}).Errorf("could not delete account")
		return err
	}

	return nil
}
