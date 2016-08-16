package transaction

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const dbKey string = "Transaction"

type Transaction struct {
	Id                 string    `datastore:"-" json:"id,omitempty" description:"Id of the transaction"`
	Name               string    `json:"name" description:"Name of payer/payee"`
	AccountId          string    `json:"accountId,omitempty" description:"Id of the account holding the transaction"`
	Date               time.Time `json:"date" description:"Date of transaction"`
	Category           string    `json:"category" description:"Category of the transaction"`
	Amount             float64   `json:"amount" description:"Amount"`
	Note               string    `json:"note" description:"Note on the transaction"`
	RelatedTransaction string    `json:"relatedTransaction,omitempty" description:"A related transaction"`
}

func Get(c context.Context, accountId string) ([]*Transaction, error) {
	return getByAncestorId(c, accountId)
}

func GetByUser(c context.Context, userId string) ([]*Transaction, error) {
	return getByAncestorId(c, userId)
}

func getByAncestorId(c context.Context, id string) ([]*Transaction, error) {
	key, err := datastore.DecodeKey(id)
	if err != nil {
		log.Errorf(c, "could not get key: %+v", err)
		return nil, err
	}
	transactions := make([]*Transaction, 0)
	q := datastore.NewQuery(dbKey).Ancestor(key).Order("-Date")
	keys, err := q.GetAll(c, &transactions)
	if err != nil {
		log.Errorf(c, "failed to fetch transactions: %+v", err)
		return nil, err
	}

	for idx, transaction := range transactions {
		transaction.Id = keys[idx].Encode()
	}

	return transactions, nil
}

func getById(c context.Context, id string) (*Transaction, error) {
	key, err := datastore.DecodeKey(id)
	if err != nil {
		log.Errorf(c, "could not get key: %+v", err)
		return nil, err
	}
	t := Transaction{}
	err = datastore.Get(c, key, &t)
	if err != nil {
		log.Errorf(c, "failed to fetch transaction: %+v", err)
		return nil, err
	}

	t.Id = id
	return &t, nil
}

func New(c context.Context, accountId string, transaction *Transaction) (*Transaction, error) {
	accountKey, err := datastore.DecodeKey(accountId)
	if err != nil {
		log.Errorf(c, "could not get account key: %+v", err)
		return nil, err
	}
	transaction.AccountId = accountId
	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, dbKey, accountKey), transaction)
	if err != nil {
		log.Errorf(c, "could not save new transaction: %+v", err)
		return nil, err
	}

	transaction.Id = key.Encode()
	return transaction, nil
}

func Update(c context.Context, transaction *Transaction, transactionId string) (*Transaction, error) {
	transactionKey, err := datastore.DecodeKey(transactionId)
	if err != nil {
		log.Errorf(c, "transaction id not valid: %+v", err)
		return nil, err
	}

	saved, err := getById(c, transactionId)
	if err != nil {
		return nil, err
	}
	transaction.AccountId = saved.AccountId

	key, err := datastore.Put(c, transactionKey, transaction)
	if err != nil {
		log.Errorf(c, "could not save new transaction: %+v", err)
		return nil, err
	}

	transaction.Id = key.Encode()
	return transaction, nil
}

func Delete(c context.Context, transactionId string) error {
	key, err := datastore.DecodeKey(transactionId)
	if err != nil {
		log.Errorf(c, "transaction id not valid: %+v", err)
		return err
	}

	err = datastore.Delete(c, key)
	if err != nil {
		log.Errorf(c, "could not delete transaction: %+v", err)
		return err
	}

	return nil
}
