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
	Title              string    `json:"title" description:"Transaction title"`
	AccountId          string    `json:"accountId,omitempty" description:"Id of the account holding the transaction"`
	Type               string    `json:"type" description:"Type of transaction"`
	Date               time.Time `json:"time" description:"Date of transaction"`
	Name               string    `json:"name" description:"Name of payer/payee"`
	Category           string    `json:"category" description:"Category of the transaction"`
	Incoming           int       `json:"incoming" description:"Amount incoming"`
	Outgoing           int       `json:"outgoing" description:"Amount outgoing"`
	Note               string    `json:"note" description:"Note on the transaction"`
	RelatedTransaction string    `json:"relatedTransaction,omitempty" description:"A related transaction"`
}

func Get(c context.Context, accountId string) ([]*Transaction, error) {
	accountKey, err := datastore.DecodeKey(accountId)
	if err != nil {
		log.Errorf(c, "could not get account key: %+v", err)
		return nil, err
	}
	transactions := make([]*Transaction, 0)
	q := datastore.NewQuery(dbKey).Ancestor(accountKey)
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
	oldTransactionKey, err := datastore.DecodeKey(transactionId)
	if err != nil {
		log.Errorf(c, "transaction id not valid: %+v", err)
		return nil, err
	}

	old := &Transaction{}
	err = datastore.Get(c, oldTransactionKey, old)
	if err != nil {
		log.Errorf(c, "could not fetch transaction: %+v", err)
		return nil, err
	}

	newTransactionKey := oldTransactionKey

	// check if parent has changed
	if transaction.AccountId != old.AccountId {
		accountKey, err := datastore.DecodeKey(transaction.AccountId)
		if err != nil {
			log.Errorf(c, "could not decode new account id: %+v", err)
			return nil, err
		}
		newTransactionKey = datastore.NewIncompleteKey(c, dbKey, accountKey)
	}

	key, err := datastore.Put(c, newTransactionKey, transaction)
	if err != nil {
		log.Errorf(c, "could not save new transaction: %+v", err)
		return nil, err
	}

	// if parent changed, clean up old transaction
	if transaction.AccountId != old.AccountId {
		err = datastore.Delete(c, oldTransactionKey)
		if err != nil {
			log.Errorf(c, "saved new transaction but could not delete old one: %+v", err)
			return nil, err
		}
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
