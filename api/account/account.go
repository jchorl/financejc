package account

import (
	"context"

	"github.com/Sirupsen/logrus"

	"github.com/jchorl/financejc/api/util"
	"github.com/jchorl/financejc/constants"
)

type Account struct {
	Id          int     `json:"id,omitempty"`
	Name        string  `json:"name"`
	Currency    string  `json:"currency"`
	User        uint    `json:"-"`
	FutureValue float64 `json:"futureValue"`
}

func Get(c context.Context) ([]*Account, error) {
	userId, err := util.UserIdFromContext(c)
	if err != nil {
		return nil, err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	accounts := make([]*Account, 0)
	rows, err := db.Query("SELECT a.id, a.name, a.currency, a.userId, COALESCE(SUM(t.amount), 0) FROM accounts a LEFT JOIN transactions t on t.account=a.id WHERE a.userId = $1 GROUP BY a.id", userId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"userId": userId,
		}).Error("failed to fetch accounts")
		return nil, err
	}

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.Id, &account.Name, &account.Currency, &account.User, &account.FutureValue); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err,
				"userId": userId,
			}).Error("failed to scan into account")
			return nil, err
		}

		accounts = append(accounts, &account)
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"userId": userId,
		}).Error("failed to get accounts from rows")
		return nil, err
	}

	return accounts, nil
}

func New(c context.Context, account *Account) (*Account, error) {
	userId, err := util.UserIdFromContext(c)
	if err != nil {
		return nil, err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	account.User = userId

	_, valid := constants.CurrencyCodeToName[account.Currency]
	if !valid {
		return nil, constants.InvalidCurrency
	}

	var id int
	err = db.QueryRow("INSERT INTO accounts(name, currency, userId) VALUES($1, $2, $3) RETURNING id", account.Name, account.Currency, account.User).Scan(&id)
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
	userId, err := util.UserIdFromContext(c)
	if err != nil {
		return nil, err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	_, valid := constants.CurrencyCodeToName[account.Currency]
	if !valid {
		return nil, constants.InvalidCurrency
	}

	_, err = db.Exec("UPDATE accounts SET name = $1, currency = $2 WHERE id = $3 AND userId = $4", account.Name, account.Currency, account.Id, userId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"account": account,
		}).Errorf("failed to update account row")
		return nil, err
	}

	return account, nil
}

func Delete(c context.Context, accountId int) error {
	userId, err := util.UserIdFromContext(c)
	if err != nil {
		return err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM accounts WHERE id = $1 AND userId = $2", accountId, userId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountId,
		}).Errorf("could not delete account")
		return err
	}

	return nil
}
