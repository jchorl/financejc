package batchTransfer

import (
	"context"
	"encoding/json"

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

func encode(data fjcData) (string, error) {
	encB, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(encB), nil
}
