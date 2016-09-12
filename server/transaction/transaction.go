package transaction

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jchorl/financejc/constants"
	"github.com/jchorl/financejc/server/util"

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
	Account            int       `json:"account" description:"The account Id that the transaction belongs to"`
}

type transactionDB struct {
	Id                 int
	Name               string
	Occurred           time.Time
	Category           sql.NullString
	Amount             float64
	Note               sql.NullString
	RelatedTransaction sql.NullInt64
	Account            int
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

func Get(c context.Context, accountId int, nextEncoded string) (Transactions, error) {
	userId := c.Value(constants.CTX_USER).(int)
	db := c.Value(constants.CTX_DB).(util.DB)

	valid, err := userOwnsAccount(c, userId, accountId)
	if err != nil {
		return Transactions{}, constants.Forbidden
	} else if !valid {
		return Transactions{}, constants.Forbidden
	}

	transactions := Transactions{}

	reference := time.Now()
	offset := 0
	if nextEncoded != "" {
		decoded, err := decodeNextPage(nextEncoded)
		if err != nil {
			return Transactions{}, err
		}

		reference, offset = decoded.Reference, decoded.Offset
	}

	rows, err := db.Query("SELECT id, name, occurred, category, amount, note, relatedTransaction, account FROM transactions WHERE account = $1 AND occurred < $2 ORDER BY occurred DESC, id LIMIT $3 OFFSET $4", accountId, reference, limitPerQuery, offset)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":      err,
			"Account ID": accountId,
			"Next":       nextEncoded,
		}).Error("failed to fetch transactions")
		return Transactions{}, err
	}

	for rows.Next() {
		var transaction transactionDB
		if err := rows.Scan(&transaction.Id, &transaction.Name, &transaction.Occurred, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.RelatedTransaction, &transaction.Account); err != nil {
			logrus.WithFields(logrus.Fields{
				"Error":      err,
				"Account ID": accountId,
				"Next":       nextEncoded,
			}).Error("failed to scan into transaction")
			return Transactions{}, err
		}

		transactions.Transactions = append(transactions.Transactions, fromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":      err,
			"Account ID": accountId,
			"Next":       nextEncoded,
		}).Error("failed to get transactions from rows")
		return Transactions{}, err
	}

	if len(transactions.Transactions) == limitPerQuery {
		next, err := encodeNextPage(nextPageParams{reference, offset + limitPerQuery})
		if err != nil {
			return Transactions{}, err
		}

		transactions.NextLink = next
	}

	return transactions, nil
}

func New(c context.Context, transaction *Transaction) (*Transaction, error) {
	userId := c.Value(constants.CTX_USER).(int)
	db := c.Value(constants.CTX_DB).(util.DB)

	valid, err := userOwnsAccount(c, userId, transaction.Account)
	if err != nil {
		return nil, constants.Forbidden
	} else if !valid {
		return nil, constants.Forbidden
	}

	tdb := toDB(*transaction)
	var id int
	err = db.QueryRow("INSERT INTO transactions(name, occurred, category, amount, note, relatedTransaction, account) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id", tdb.Name, tdb.Occurred, tdb.Category, tdb.Amount, tdb.Note, tdb.RelatedTransaction, tdb.Account).Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":         err,
			"TransactionDB": tdb,
			"Transaction":   transaction,
		}).Errorf("failed to insert transaction row")
		return nil, err
	}

	transaction.Id = id
	return transaction, nil
}

func Update(c context.Context, transaction *Transaction) (*Transaction, error) {
	userId := c.Value(constants.CTX_USER).(int)
	db := c.Value(constants.CTX_DB).(util.DB)

	valid, err := userOwnsAccount(c, userId, transaction.Account)
	if err != nil {
		return nil, constants.Forbidden
	} else if !valid {
		return nil, constants.Forbidden
	}

	tdb := toDB(*transaction)
	_, err = db.Exec("UPDATE transactions SET name = $1, occurred = $2, category = $3, amount = $4, note = $5, relatedTransaction = $6 WHERE id = $7", tdb.Name, tdb.Occurred, tdb.Category, tdb.Amount, tdb.Note, tdb.RelatedTransaction, tdb.Id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":         err,
			"TransactionDB": tdb,
			"Transaction":   transaction,
		}).Errorf("failed to update transaction row")
		return nil, err
	}

	return transaction, nil
}

func Delete(c context.Context, transactionId int) error {
	userId := c.Value(constants.CTX_USER).(int)
	db := c.Value(constants.CTX_DB).(util.DB)

	valid, err := userOwnsTransaction(c, userId, transactionId)
	if err != nil {
		return constants.Forbidden
	} else if !valid {
		return constants.Forbidden
	}

	_, err = db.Exec("DELETE FROM transactions WHERE id = $1", transactionId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":          err,
			"Transaction ID": transactionId,
		}).Errorf("could not delete transaction")
		return err
	}

	return nil
}

func userOwnsAccount(c context.Context, user int, account int) (bool, error) {
	db := c.Value(constants.CTX_DB).(util.DB)
	var owner int
	err := db.QueryRow("SELECT userId FROM accounts WHERE id = $1", account).Scan(&owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":   err,
			"User":    user,
			"Account": account,
		}).Error("error checking owner of account")
		return false, err
	}

	return owner == user, nil
}

func userOwnsTransaction(c context.Context, user int, transaction int) (bool, error) {
	db := c.Value(constants.CTX_DB).(util.DB)
	var owner int
	err := db.QueryRow("SELECT a.userId FROM accounts a JOIN transactions t ON t.account = a.id WHERE t.id = $1", transaction).Scan(&owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":       err,
			"User":        user,
			"Transaction": transaction,
		}).Error("error checking owner of transaction")
		return false, err
	}

	return owner == user, nil
}

func encodeNextPage(decoded nextPageParams) (string, error) {
	bts, err := json.Marshal(decoded)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":   err,
			"Decoded": decoded,
		}).Error("could not encode next page parameter")
		return "", err
	}

	return string(bts), nil
}

func decodeNextPage(encoded string) (nextPageParams, error) {
	var decoded nextPageParams
	err := json.Unmarshal([]byte(encoded), &decoded)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":   err,
			"Encoded": encoded,
		}).Error("could not decode next page parameter")
		return nextPageParams{}, err
	}

	return decoded, nil
}

func toDB(transaction Transaction) *transactionDB {
	return &transactionDB{
		Id:                 transaction.Id,
		Name:               transaction.Name,
		Occurred:           transaction.Date,
		Category:           util.ToNullString(transaction.Category),
		Amount:             transaction.Amount,
		Note:               util.ToNullString(transaction.Note),
		RelatedTransaction: util.ToNullInt(transaction.RelatedTransaction),
		Account:            transaction.Account,
	}
}

func fromDB(transaction transactionDB) *Transaction {
	return &Transaction{
		Id:                 transaction.Id,
		Name:               transaction.Name,
		Date:               transaction.Occurred,
		Category:           util.FromNullString(transaction.Category),
		Amount:             transaction.Amount,
		Note:               util.FromNullString(transaction.Note),
		RelatedTransaction: util.FromNullInt(transaction.RelatedTransaction),
		Account:            transaction.Account,
	}
}
