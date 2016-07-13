package account

import (
	"errors"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"github.com/jchorl/financejc/currency"
	"github.com/jchorl/financejc/transaction"
)

const dbKey string = "Account"

var InvalidCurrency = errors.New("Currency is invalid")

type Account struct {
	Id            string                     `datastore:"-" json:"id,omitempty" description:"Id of the account"`
	Name          string                     `json:"name" description:"Name of the account"`
	Currency      string                     `json:"currency" description:"Currency for the account"`
	Balance       float64                    `json:"balance" description:"Current balance" datastore:"-"`
	FutureBalance float64                    `json:"futureBalance" description:"Balance including future transactions" datastore:"-"`
	Transactions  []*transaction.Transaction `json:"transactions" datastore:"-" description:"Transactions for the account"`
}

func Get(c context.Context, userId string) ([]*Account, error) {
	userKey, err := datastore.DecodeKey(userId)
	if err != nil {
		log.Errorf(c, "could not get user key: %+v", err)
		return nil, err
	}
	accounts := make([]*Account, 0)
	q := datastore.NewQuery(dbKey).Ancestor(userKey)
	keys, err := q.GetAll(c, &accounts)
	if err != nil {
		log.Errorf(c, "failed to fetch accounts: %+v", err)
		return nil, err
	}

	for idx, acc := range accounts {
		acc.Id = keys[idx].Encode()
	}

	err = setBalancesAndTransactions(c, userId, accounts)
	if err != nil {
		log.Errorf(c, "failed to set balances and transactions: %+v", err)
		return nil, err
	}

	return accounts, nil
}

func setBalancesAndTransactions(c context.Context, userId string, accounts []*Account) error {
	keyToAccount := make(map[string]*Account)
	for _, acc := range accounts {
		keyToAccount[acc.Id] = acc
		acc.Balance = 0
		acc.FutureBalance = 0
		acc.Transactions = make([]*transaction.Transaction, 0)
	}

	transactions, err := transaction.GetByUser(c, userId)
	if err != nil {
		log.Errorf(c, "failed to fetch transactions: %+v", err)
		return err
	}

	now := time.Now()

	var account *Account

	for _, tr := range transactions {
		account = keyToAccount[tr.AccountId]
		account.Transactions = append(account.Transactions, tr)
		account.FutureBalance += tr.Amount
		if tr.Date.Before(now) {
			account.Balance += tr.Amount
		}
	}

	return nil
}

func New(c context.Context, userId string, account *Account) (*Account, error) {
	return putAccount(c, userId, account, true, "")
}

func Update(c context.Context, account *Account, accountId string) (*Account, error) {
	return putAccount(c, "", account, false, accountId)
}

func Delete(c context.Context, accountId string) error {
	key, err := datastore.DecodeKey(accountId)
	if err != nil {
		log.Errorf(c, "account id not valid: %+v", err)
		return err
	}

	err = datastore.Delete(c, key)
	if err != nil {
		log.Errorf(c, "could not delete account: %+v", err)
		return err
	}

	return nil
}

func putAccount(c context.Context, userId string, account *Account, newAccount bool, accountId string) (*Account, error) {
	_, valid := currency.CodeToName[account.Currency]
	if !valid {
		return nil, InvalidCurrency
	}

	var accountKey *datastore.Key
	if newAccount {
		userKey, err := datastore.DecodeKey(userId)
		if err != nil {
			log.Errorf(c, "could not get user key: %+v", err)
			return nil, err
		}
		accountKey = datastore.NewIncompleteKey(c, dbKey, userKey)
	} else {
		var err error
		accountKey, err = datastore.DecodeKey(accountId)
		if err != nil {
			log.Errorf(c, "account id not valid: %+v", err)
			return nil, err
		}
	}

	key, err := datastore.Put(c, accountKey, account)
	if err != nil {
		log.Errorf(c, "could not save new account: %+v", err)
		return nil, err
	}

	account.Id = key.Encode()
	account.Transactions = make([]*transaction.Transaction, 0)
	return account, nil
}
