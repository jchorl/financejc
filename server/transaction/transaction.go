package transaction

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jchorl/financejc/constants"

	"github.com/Sirupsen/logrus"
)

const (
	limitPerQuery = 25
)

type Transactions struct {
	NextLink     string
	Transactions []*Transaction `json:"transactions" description:"The transactions returned"`
}

type Transaction struct {
	Id                 int       `json:"id,omitempty" description:"Id of the transaction"`
	Name               string    `json:"name" description:"Name of payer/payee"`
	Date               time.Time `json:"date" description:"Date of transaction"`
	Category           string    `json:"category" description:"Category of the transaction"`
	Amount             float64   `json:"amount" description:"Amount"`
	Note               string    `json:"note" description:"Note on the transaction"`
	RelatedTransaction int       `json:"relatedTransaction,omitempty" description:"A related transaction"`
	Account            int       `json:"-"`
}

type nextPageParams struct {
	Reference time.Time
	Offset    int
}

func (t Transactions) Next() string {
	return t.NextLink
}

func (t Transactions) Values() (ret []interface{}) {
	for _, tr := range t.Transactions {
		ret = append(ret, tr)
	}

	return ret
}

func Get(c context.Context, accountId, nextEncoded string) (Transactions, error) {
	db := c.Value(constants.CTX_DB).(sql.DB)
	transactions := Transactions{
		Transactions: make([]*Transaction, 0),
	}

	reference := time.Now()
	offset := 0
	if nextEncoded != "" {
		decoded, err := decodeNextPage(nextEncoded)
		if err != nil {
			return Transactions{}, err
		}

		reference, offset = decoded.Reference, decoded.Offset
	}

	rows, err := db.Query("SELECT * FROM transactions WHERE account = $1 AND date < $2 ORDER BY date DESC, id LIMIT $3 OFFSET $4", accountId, reference, limitPerQuery, offset)
	if err != nil {
		logrus.WithField("Error", err).Error("failed to fetch transactions")
		return Transactions{}, err
	}

	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction); err != nil {
			logrus.WithField("Error", err).Error("failed to scan into transaction")
			return Transactions{}, err
		}

		transactions.Transactions = append(transactions.Transactions, &transaction)
	}
	if err := rows.Err(); err != nil {
		logrus.WithField("Error", err).Error("failed to get transactions from rows")
		return Transactions{}, err
	}

	return transactions, nil
}

func New(c context.Context, transaction *Transaction) (*Transaction, error) {
	db := c.Value(constants.CTX_DB).(sql.DB)
	var id int
	err := db.QueryRow("INSERT INTO transactions(name, date, category, amount, note, relatedTransaction, account) VALUES($1, $2, $3, $4, $5, $6, %7) RETURNING id", transaction.Name, transaction.Date, transaction.Category, transaction.Amount, transaction.Note, transaction.RelatedTransaction, transaction.Account).Scan(&id)
	if err != nil {
		logrus.WithField("Error", err).Errorf("failed to insert transaction row")
		return nil, err
	}

	transaction.Id = id
	return transaction, nil
}

func Update(c context.Context, transaction *Transaction) (*Transaction, error) {
	db := c.Value(constants.CTX_DB).(sql.DB)
	_, err := db.Exec("UPDATE transactions SET name = $1, date = $2, category = $3, amount = $4, note = $5, relatedTransaction = $6, account = $7 WHERE id = $8", transaction.Name, transaction.Date, transaction.Category, transaction.Amount, transaction.Note, transaction.RelatedTransaction, transaction.Account, transaction.Id)
	if err != nil {
		logrus.WithField("Error", err).Errorf("failed to update transaction row")
		return nil, err
	}

	return transaction, nil
}

func Delete(c context.Context, transactionId int) error {
	db := c.Value(constants.CTX_DB).(sql.DB)

	_, err := db.Exec("DELETE FROM transactions WHERE id = $1", transactionId)
	if err != nil {
		logrus.WithField("Error", err).Errorf("could not delete transaction")
		return err
	}

	return nil
}

func encodeNextPage(decoded nextPageParams) (string, error) {
	bts, err := json.Marshal(decoded)
	if err != nil {
		logrus.WithField("Error", err).Error("could not encode next page parameter")
		return "", err
	}

	return string(bts), nil
}

func decodeNextPage(encoded string) (nextPageParams, error) {
	var decoded nextPageParams
	err := json.Unmarshal([]byte(encoded), &decoded)
	if err != nil {
		logrus.WithField("Error", err).Error("could not decode next page parameter")
		return nextPageParams{}, err
	}

	return decoded, nil
}
