package transaction

import (
	"context"
	"database/sql"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/jchorl/financejc/api/util"
	"github.com/jchorl/financejc/constants"
)

type RecurringTransaction struct {
	Id                  int         `json:"id,omitempty"`
	Transaction         Transaction `json:"transaction"`
	ScheduleType        string      `json:"scheduleType"`
	SecondsBetween      *int        `json:"secondsBetween"`
	DayOf               *int        `json:"dayOf"`
	SecondsBeforeToPost int         `json:"secondsBeforeToPost"`
}

type recurringTransactionDB struct {
	Id         int
	Name       string
	NextOccurs time.Time
	Category   string
	Amount     int
	Note       string
	AccountId  int

	ScheduleType        string
	SecondsBetween      sql.NullInt64
	DayOf               sql.NullInt64
	SecondsBeforeToPost int
}

func GenRecurringTransactions() {
	// query for all recurring transactions where the next occurrance is within the time period before to post it to the account
	// "select id, name, nextOccurs, category, amount, note, account, scheduleType, secondsBetween, dayOf, secondsBeforeToPost from recurringTransactions where nextOccurs - interval '1 second' * secondsBeforeToPost <= currentTimestamp;"
	// for each one, generate a transaction
	// calculate when the transaction should next run
	// update the recurring transaction
}

func GetRecurring(c context.Context, accountId int) ([]RecurringTransaction, error) {
	transactions := []RecurringTransaction{}
	db, err := util.DBFromContext(c)
	if err != nil {
		return transactions, err
	}

	valid, err := userOwnsAccount(c, accountId)
	if err != nil || !valid {
		return transactions, constants.Forbidden
	}

	rows, err := db.Query("SELECT id, name, nextOccurs, category, amount, note, account, scheduleType, secondsBetween, dayOf, secondsBeforeToPost FROM recurringTransactions WHERE account = $1", accountId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountId,
		}).Error("failed to fetch recurring transactions")
		return transactions, err
	}

	for rows.Next() {
		var transaction recurringTransactionDB
		if err := rows.Scan(&transaction.Id, &transaction.Name, &transaction.NextOccurs, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.AccountId, &transaction.ScheduleType, &transaction.SecondsBetween, &transaction.DayOf, &transaction.SecondsBeforeToPost); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"accountId": accountId,
			}).Error("failed to scan into recurring transaction")
			return transactions, err
		}

		transactions = append(transactions, recurringFromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountId,
		}).Error("failed to get recurring transactions from rows")
		return transactions, err
	}

	return transactions, nil
}

func NewRecurring(c context.Context, transaction *RecurringTransaction) (*RecurringTransaction, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	valid, err := userOwnsAccount(c, transaction.Transaction.AccountId)
	if err != nil || !valid {
		return nil, constants.Forbidden
	}

	tdb := recurringToDB(*transaction)
	var id int
	err = db.QueryRow("INSERT INTO recurringTransactions(name, nextOccurs, category, amount, note, account, scheduleType, secondsBetween, dayOf, secondsBeforeToPost) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id", tdb.Name, tdb.NextOccurs, tdb.Category, tdb.Amount, tdb.Note, tdb.AccountId, tdb.ScheduleType, tdb.SecondsBetween, tdb.DayOf, tdb.SecondsBeforeToPost).Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":                  err,
			"recurringTransactionDB": tdb,
			"recurringTransaction":   transaction,
		}).Errorf("failed to insert recurring transaction row")
		return nil, err
	}

	transaction.Id = id
	return transaction, nil
}

func recurringToDB(transaction RecurringTransaction) *recurringTransactionDB {
	return &recurringTransactionDB{
		Id:         transaction.Id,
		Name:       transaction.Transaction.Name,
		NextOccurs: transaction.Transaction.Date,
		Category:   transaction.Transaction.Category,
		Amount:     transaction.Transaction.Amount,
		Note:       transaction.Transaction.Note,
		AccountId:  transaction.Transaction.AccountId,

		ScheduleType:        transaction.ScheduleType,
		SecondsBetween:      util.ToNullInt(transaction.SecondsBetween),
		DayOf:               util.ToNullInt(transaction.DayOf),
		SecondsBeforeToPost: transaction.SecondsBeforeToPost,
	}
}

func recurringFromDB(transaction recurringTransactionDB) RecurringTransaction {
	return RecurringTransaction{
		Transaction: Transaction{
			Id:        transaction.Id,
			Name:      transaction.Name,
			Date:      transaction.NextOccurs,
			Category:  transaction.Category,
			Amount:    transaction.Amount,
			Note:      transaction.Note,
			AccountId: transaction.AccountId,
		},

		ScheduleType:        transaction.ScheduleType,
		SecondsBetween:      util.FromNullInt(transaction.SecondsBetween),
		DayOf:               util.FromNullInt(transaction.DayOf),
		SecondsBeforeToPost: transaction.SecondsBeforeToPost,
	}
}
