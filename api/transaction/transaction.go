package transaction

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/jchorl/financejc/api/util"
	"github.com/jchorl/financejc/constants"
)

const (
	limitPerQuery = 25
)

type Transactions struct {
	NextLink     string
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Id                   int       `json:"id,omitempty"`
	Name                 string    `json:"name"`
	Date                 time.Time `json:"date"`
	Category             string    `json:"category"`
	Amount               int       `json:"amount"`
	Note                 string    `json:"note"`
	RelatedTransactionId int       `json:"relatedTransactionId,omitempty"`
	AccountId            int       `json:"accountId"`
}

type transactionDB struct {
	Id                   int
	Name                 string
	Occurred             time.Time
	Category             sql.NullString
	Amount               int
	Note                 sql.NullString
	RelatedTransactionId sql.NullInt64
	AccountId            int
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
	db, err := util.DBFromContext(c)
	if err != nil {
		return Transactions{}, err
	}

	valid, err := userOwnsAccount(c, accountId)
	if err != nil || !valid {
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

	rows, err := db.Query("SELECT id, name, occurred, category, amount, note, relatedTransactionId, accountId FROM transactions WHERE accountId = $1 AND occurred < $2 ORDER BY occurred DESC, id LIMIT $3 OFFSET $4", accountId, reference, limitPerQuery, offset)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountId,
			"next":      nextEncoded,
		}).Error("failed to fetch transactions")
		return Transactions{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction transactionDB
		if err := rows.Scan(&transaction.Id, &transaction.Name, &transaction.Occurred, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.RelatedTransactionId, &transaction.AccountId); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"accountId": accountId,
				"next":      nextEncoded,
			}).Error("failed to scan into transaction")
			return Transactions{}, err
		}

		transactions.Transactions = append(transactions.Transactions, fromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountId,
			"next":      nextEncoded,
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

func GetFuture(c context.Context, accountId int, reference *time.Time) ([]Transaction, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	valid, err := userOwnsAccount(c, accountId)
	if err != nil || !valid {
		return nil, constants.Forbidden
	}

	transactions := []Transaction{}

	if reference == nil {
		now := time.Now()
		reference = &now
	}
	rows, err := db.Query("SELECT id, name, occurred, category, amount, note, relatedTransactionId, accountId FROM transactions WHERE accountId = $1 AND occurred > $2 ORDER BY occurred DESC, id", accountId, reference)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountId,
		}).Error("failed to fetch future transactions")
		return transactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction transactionDB
		if err := rows.Scan(&transaction.Id, &transaction.Name, &transaction.Occurred, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.RelatedTransactionId, &transaction.AccountId); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"accountId": accountId,
			}).Error("failed to scan into transaction for future fetch")
			return transactions, err
		}

		transactions = append(transactions, fromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountId,
		}).Error("failed to get transactions from rows for future fetch")
		return transactions, err
	}

	return transactions, nil
}

func New(c context.Context, transaction *Transaction) (*Transaction, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	valid, err := userOwnsAccount(c, transaction.AccountId)
	if err != nil || !valid {
		return nil, constants.Forbidden
	}

	tdb := toDB(*transaction)
	var id int
	err = db.QueryRow("INSERT INTO transactions(name, occurred, category, amount, note, relatedTransactionId, accountId) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id", tdb.Name, tdb.Occurred, tdb.Category, tdb.Amount, tdb.Note, tdb.RelatedTransactionId, tdb.AccountId).Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":         err,
			"transactionDB": tdb,
			"transaction":   transaction,
		}).Errorf("failed to insert transaction row")
		return nil, err
	}

	transaction.Id = id
	return transaction, nil
}

func Update(c context.Context, transaction *Transaction) (*Transaction, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	valid, err := userOwnsAccount(c, transaction.AccountId)
	if err != nil || !valid {
		return nil, constants.Forbidden
	}

	tdb := toDB(*transaction)
	_, err = db.Exec("UPDATE transactions SET name = $1, occurred = $2, category = $3, amount = $4, note = $5, relatedTransactionId = $6 WHERE id = $7", tdb.Name, tdb.Occurred, tdb.Category, tdb.Amount, tdb.Note, tdb.RelatedTransactionId, tdb.Id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":         err,
			"transactionDB": tdb,
			"transaction":   transaction,
		}).Errorf("failed to update transaction row")
		return nil, err
	}

	return transaction, nil
}

func Delete(c context.Context, transactionId int) error {
	db, err := util.DBFromContext(c)
	if err != nil {
		return err
	}

	valid, err := userOwnsTransaction(c, transactionId)
	if err != nil {
		return constants.Forbidden
	} else if !valid {
		return constants.Forbidden
	}

	_, err = db.Exec("DELETE FROM transactions WHERE id = $1", transactionId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":         err,
			"transactionID": transactionId,
		}).Errorf("could not delete transaction")
		return err
	}

	return nil
}

func userOwnsAccount(c context.Context, account int) (bool, error) {
	userId, err := util.UserIdFromContext(c)
	if err != nil {
		return false, err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return false, err
	}

	var owner uint
	err = db.QueryRow("SELECT userId FROM accounts WHERE id = $1", account).Scan(&owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"userId":  userId,
			"account": account,
		}).Error("error checking owner of account")
		return false, err
	}

	return owner == userId, nil
}

func userOwnsTransaction(c context.Context, transaction int) (bool, error) {
	userId, err := util.UserIdFromContext(c)
	if err != nil {
		return false, err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return false, err
	}

	var owner uint
	err = db.QueryRow("SELECT a.userId FROM accounts a JOIN transactions t ON t.accountId = a.id WHERE t.id = $1", transaction).Scan(&owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":       err,
			"userId":      userId,
			"transaction": transaction,
		}).Error("error checking owner of transaction")
		return false, err
	}

	return owner == userId, nil
}

func encodeNextPage(decoded nextPageParams) (string, error) {
	bts, err := json.Marshal(decoded)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"decoded": decoded,
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
			"error":   err,
			"encoded": encoded,
		}).Error("could not decode next page parameter")
		return nextPageParams{}, err
	}

	return decoded, nil
}

func toDB(transaction Transaction) *transactionDB {
	return &transactionDB{
		Id:                   transaction.Id,
		Name:                 transaction.Name,
		Occurred:             transaction.Date,
		Category:             util.ToNullStringNonEmpty(transaction.Category),
		Amount:               transaction.Amount,
		Note:                 util.ToNullStringNonEmpty(transaction.Note),
		RelatedTransactionId: util.ToNullIntNonZero(transaction.RelatedTransactionId),
		AccountId:            transaction.AccountId,
	}
}

func fromDB(transaction transactionDB) Transaction {
	return Transaction{
		Id:                   transaction.Id,
		Name:                 transaction.Name,
		Date:                 transaction.Occurred,
		Category:             util.FromNullStringNonEmpty(transaction.Category),
		Amount:               transaction.Amount,
		Note:                 util.FromNullStringNonEmpty(transaction.Note),
		RelatedTransactionId: util.FromNullIntNonZero(transaction.RelatedTransactionId),
		AccountId:            transaction.AccountId,
	}
}
