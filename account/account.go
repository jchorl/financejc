package account

import (
	"errors"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"currency"
	"user"
)

const dbKey string = "Account"

var InvalidCurrency = errors.New("Currency is invalid")

type Account struct {
	Id                 string    `datastore:"-" json:"id,omitempty" description:"Id of the account"`
	Name               string    `json:"name" description:"Name of the account"`
	Currency           string    `json:"currency" description:"Currency for the account"`
	Balance            int       `json:"balance" description:"Current balance"`
	FutureBalance      int       `json:"futureBalance" description:"Balance including future transactions"`
	BalancesCalculated time.Time `json:"-"`
}

func Get(c context.Context, userId string) ([]*Account, error) {
	userKey, err := user.Key(userId)
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

	return accounts, nil
}

func New(c context.Context, userId string, account *Account) (*Account, error) {
	return putAccount(c, userId, account, true, "")
}

func Update(c context.Context, userId string, account *Account, accountId string) (*Account, error) {
	return putAccount(c, userId, account, false, accountId)
}

func putAccount(c context.Context, userId string, account *Account, newAccount bool, accountId string) (*Account, error) {
	userKey, err := user.Key(userId)
	if err != nil {
		log.Errorf(c, "could not get user key: %+v", err)
		return nil, err
	}

	_, valid := currency.CodeToName[account.Currency]
	if !valid {
		return nil, InvalidCurrency
	}

	var accountKey *datastore.Key
	if newAccount {
		accountKey = datastore.NewIncompleteKey(c, dbKey, userKey)
		account.Balance = 0
		account.FutureBalance = 0
		account.BalancesCalculated = time.Now()
	} else {
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
	return account, nil
}
