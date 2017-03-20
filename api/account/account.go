package account

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/lib/pq"

	"github.com/jchorl/financejc/api/util"
	"github.com/jchorl/financejc/constants"
)

// Account is a user's bank account
type Account struct {
	ID          int     `json:"id,omitempty"`
	Name        string  `json:"name"`
	Currency    string  `json:"currency"`
	User        uint    `json:"user"`
	FutureValue float64 `json:"futureValue"`
}

// Get fetches all accounts of a user
func Get(c context.Context) ([]*Account, error) {
	userID, err := util.UserIDFromContext(c)
	if err != nil {
		return nil, err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	rows, err := db.Query("SELECT a.id, a.name, a.currency, a.user_id, COALESCE(SUM(t.amount), 0) FROM accounts a LEFT JOIN transactions t on t.account_id=a.id WHERE a.user_id = $1 GROUP BY a.id", userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"userId": userID,
		}).Error("failed to fetch accounts")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Name, &account.Currency, &account.User, &account.FutureValue); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err,
				"userId": userID,
			}).Error("failed to scan into account")
			return nil, err
		}

		accounts = append(accounts, &account)
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"userId": userID,
		}).Error("failed to get accounts from rows")
		return nil, err
	}

	return accounts, nil
}

// BatchImport batch imports accounts
func BatchImport(c context.Context, accounts []Account) error {
	userID, err := util.UserIDFromContext(c)
	if err != nil || !util.IsUserAdmin(userID) {
		return constants.ErrForbidden
	}

	db, err := util.SQLDBFromContext(c)
	if err != nil {
		return err
	}

	txn, err := db.Begin()
	if err != nil {
		logrus.WithError(err).Error("unable to begin transaction when batch inserting accounts")
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn("accounts", "id", "name", "currency", "user_id"))
	if err != nil {
		logrus.WithError(err).Error("unable to begin copy in when batch inserting accounts")
		return err
	}

	for _, account := range accounts {
		_, err = stmt.Exec(account.ID, account.Name, account.Currency, account.User)
		if err != nil {
			logrus.WithError(err).Error("unable to exec transaction copy when batch inserting accounts")
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		logrus.WithError(err).Error("unable to exec batch account copy when batch inserting accounts")
		return err
	}

	err = stmt.Close()
	if err != nil {
		logrus.WithError(err).Error("unable to close account copy when batch inserting accounts")
		return err
	}

	err = txn.Commit()
	if err != nil {
		logrus.WithError(err).Error("unable to commit account copy when batch inserting accounts")
		return err
	}

	return nil
}

// GetAll queries for all accounts
func GetAll(c context.Context) ([]Account, error) {
	userID, err := util.UserIDFromContext(c)
	if err != nil || !util.IsUserAdmin(userID) {
		return nil, constants.ErrForbidden
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	accounts := []Account{}
	rows, err := db.Query("SELECT id, name, currency, user_id FROM accounts")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to fetch all accounts")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Name, &account.Currency, &account.User); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("failed to scan into account")
			return nil, err
		}

		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to get all accounts from rows")
		return nil, err
	}

	return accounts, nil
}

// New creates a new account
func New(c context.Context, account *Account) (*Account, error) {
	userID, err := util.UserIDFromContext(c)
	if err != nil {
		return nil, err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	account.User = userID

	_, valid := constants.CurrencyInfo[account.Currency]
	if !valid {
		return nil, constants.ErrInvalidCurrency
	}

	var id int
	err = db.QueryRow("INSERT INTO accounts(name, currency, user_id) VALUES($1, $2, $3) RETURNING id", account.Name, account.Currency, account.User).Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":   err,
			"Account": account,
		}).Errorf("failed to insert account row")
		return nil, err
	}

	account.ID = id
	return account, nil
}

// Update updates an account
func Update(c context.Context, account *Account) (*Account, error) {
	valid, err := util.UserOwnsAccount(c, account.ID)
	if err != nil || !valid {
		return nil, constants.ErrForbidden
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	_, valid = constants.CurrencyInfo[account.Currency]
	if !valid {
		return nil, constants.ErrInvalidCurrency
	}

	_, err = db.Exec("UPDATE accounts SET name = $1, currency = $2 WHERE id = $3", account.Name, account.Currency, account.ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"account": account,
		}).Errorf("failed to update account row")
		return nil, err
	}

	return account, nil
}

// Delete deletes an account
func Delete(c context.Context, accountID int) error {
	valid, err := util.UserOwnsAccount(c, accountID)
	if err != nil || !valid {
		return constants.ErrForbidden
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM accounts WHERE id = $1", accountID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountID,
		}).Errorf("could not delete account")
		return err
	}

	return nil
}
