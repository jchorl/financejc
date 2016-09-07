package transaction

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const (
	dbKey         = "Transaction"
	limitPerQuery = 25
)

type Transactions struct {
	Next         string
	Transactions []*Transaction `json:"transactions" description:"The transactions returned"`
}

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

func Get(c context.Context, accountId, cursorEncoded string) (Transactions, error) {
	accountKey, err := datastore.DecodeKey(accountId)
	if err != nil {
		log.Errorf(c, "could not get account key: %+v", err)
		return Transactions{}, err
	}

	transactions := Transactions{
		Transactions: make([]*Transaction, 0),
	}

	q := datastore.NewQuery(dbKey).Ancestor(accountKey).Order("-Date")

	if cursorEncoded != "" {
		cursor, err := datastore.DecodeCursor(cursorEncoded)
		if err != nil {
			log.Errorf(c, "could not get key: %+v", err)
			return Transactions{}, err
		}

		q = q.Start(cursor)
	}

	for i, it := 0, q.Run(c); i < limitPerQuery; i++ {
		var transaction Transaction
		key, err := it.Next(&transaction)
		if err == datastore.Done {
			break
		}
		if err != nil {
			log.Errorf(c, "could not get next transaction: %+v", err)
			return Transactions{}, err
		}
		transaction.Id = key.Encode()
		transactions.Transactions = append(transactions.Transactions, &transaction)

		if i+1 == limitPerQuery {
			cursor, err := it.Cursor()
			if err != nil {
				log.Errorf(c, "could not get cursor from iterator: %+v", err)
				return Transactions{}, err
			}
			transactions.Next = cursor.String()
		}
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
