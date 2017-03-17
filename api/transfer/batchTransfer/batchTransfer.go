package batchTransfer

import (
	"context"
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/jchorl/financejc/api/account"
	"github.com/jchorl/financejc/api/transaction"
	"github.com/jchorl/financejc/api/user"
	"github.com/jchorl/financejc/api/util"
	"github.com/jchorl/financejc/constants"
)

type fjcData struct {
	Users                 []user.User                        `json:"users"`
	Accounts              []account.Account                  `json:"accounts"`
	Transactions          []transaction.Transaction          `json:"transactions"`
	RecurringTransactions []transaction.RecurringTransaction `json:"recurringTransactions"`
	Templates             []transaction.Template             `json:"templates"`
}

// Export queries for all data, packages it up and exports it
func Export(c context.Context) (string, error) {
	userID, err := util.UserIDFromContext(c)
	if err != nil || !util.IsUserAdmin(userID) {
		return "", constants.ErrForbidden
	}

	allData := fjcData{}
	users, err := user.GetAll(c)
	if err != nil {
		return "", err
	}
	allData.Users = users

	accounts, err := account.GetAll(c)
	if err != nil {
		return "", err
	}
	allData.Accounts = accounts

	transactions, err := transaction.GetAll(c)
	if err != nil {
		return "", err
	}
	allData.Transactions = transactions

	templates, err := transaction.GetAllTemplates(c)
	if err != nil {
		return "", err
	}
	allData.Templates = templates

	recurringTransactions, err := transaction.GetAllRecurring(c)
	if err != nil {
		return "", err
	}
	allData.RecurringTransactions = recurringTransactions

	return encode(allData)
}

// Import batch imports the result of an export
func Import(c context.Context, encoded string) error {
	allData, err := decode(encoded)
	if err != nil {
		return err
	}

	if err = user.BatchImport(c, allData.Users); err != nil {
		return err
	}

	if err = account.BatchImport(c, allData.Accounts); err != nil {
		return err
	}

	if err = transaction.BatchImport(c, allData.Transactions); err != nil {
		return err
	}

	if err = transaction.BatchImportTemplates(c, allData.Templates); err != nil {
		return err
	}

	if err = transaction.BatchImportRecurringTransactions(c, allData.RecurringTransactions); err != nil {
		return err
	}

	return nil
}

func encode(data fjcData) (string, error) {
	encB, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Error("error encoding all fjc data")
		return "", err
	}

	return string(encB), nil
}

func decode(str string) (fjcData, error) {
	data := fjcData{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		logrus.WithError(err).Error("error decoding all fjc data")
		return data, err
	}

	return data, nil
}
