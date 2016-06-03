package upload

import (
	"os"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"account"
	"transaction"
)

type ParsedTransaction interface {
	GetName() (string, error)
	GetAccount() (string, error)
	GetDate() (time.Time, error)
	GetCategory() (string, error)
	GetIncoming() (float64, error)
	GetOutgoing() (float64, error)
	GetNote() (string, error)
	GetCurrency() (string, error)
}

func Upload(c context.Context, userId string, file *os.File, kind string) error {
	var parsedTransactions []ParsedTransaction
	var err error
	switch kind {
	case "TSV":
		parsedTransactions, err = TSVUpload(file)
	}
	if err != nil {
		return err
	}

	log.Debugf(c, "Finished parsing transactions")
	accounts := make(map[string]*account.Account)

	return datastore.RunInTransaction(c, func(c context.Context) error {
		for _, ptr := range parsedTransactions {
			// if the account does not exist, add it
			accountName, err := ptr.GetAccount()
			if err != nil {
				return err
			}

			act, present := accounts[accountName]
			if !present {
				curr, err := ptr.GetCurrency()
				if err != nil {
					return err
				}

				act = &account.Account{
					Name:     accountName,
					Currency: curr,
				}
				act, err = account.New(c, userId, act)
				if err != nil {
					return err
				}

				log.Debugf(c, "Saved account: %s", accountName)

				accounts[act.Name] = act
			}

			name, err := ptr.GetName()
			if err != nil {
				return err
			}
			date, err := ptr.GetDate()
			if err != nil {
				return err
			}
			category, err := ptr.GetCategory()
			if err != nil {
				return err
			}
			inc, err := ptr.GetIncoming()
			if err != nil {
				return err
			}
			out, err := ptr.GetOutgoing()
			if err != nil {
				return err
			}
			note, err := ptr.GetNote()
			if err != nil {
				return err
			}

			tr := &transaction.Transaction{
				Name:     name,
				Date:     date,
				Category: category,
				Incoming: inc,
				Outgoing: out,
				Note:     note,
			}
			_, err = transaction.New(c, act.Id, tr)
			if err != nil {
				return err
			}
			log.Debugf(c, "Saved transaction: %s", name)
		}
		return nil
	}, nil)
}
