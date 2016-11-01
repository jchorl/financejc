package transaction

import (
	"context"
	"database/sql"
	"fmt"
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

func GenRecurringTransactions(ctx context.Context) error {
	db, err := util.SQLDBFromContext(ctx)
	if err != nil {
		return err
	}

	// query for all recurring transactions where the next occurrance is within the time period before to post it to the account
	rows, err := db.Query("SELECT id, name, nextOccurs, category, amount, note, accountId, scheduleType, secondsBetween, dayOf, secondsBeforeToPost FROM recurringTransactions WHERE nextOccurs - interval '1 second' * secondsBeforeToPost <= currentTimestamp")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to fetch recurring transactions to post to accounts")
		return err
	}

	recurringTransactions := []RecurringTransaction{}
	for rows.Next() {
		var transaction recurringTransactionDB
		if err := rows.Scan(&transaction.Id, &transaction.Name, &transaction.NextOccurs, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.AccountId, &transaction.ScheduleType, &transaction.SecondsBetween, &transaction.DayOf, &transaction.SecondsBeforeToPost); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("failed to scan into recurring transaction to generate transaction")
			return err
		}

		recurringTransactions = append(recurringTransactions, recurringFromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to get recurring transactions from rows to generate transactions")
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		logrus.WithError(err).Error("could not begin transaction to generate recurring transactions")
		return err
	}

	// replace the sql Db in the context with the sql Tx
	ctx = context.WithValue(ctx, constants.CTX_DB, tx)

	for _, recurringTransaction := range recurringTransactions {
		// get the transaction from the recurring transaction template
		if _, err := New(ctx, &recurringTransaction.Transaction); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":                err,
				"recurringTransaction": recurringTransaction,
			}).Error("error adding a new transaction for a recurring transaction")
			return err
		}

		// calculate when the transaction should next run
		recurringTransaction.Transaction.Date, err = getNextRun(recurringTransaction)
		if err != nil {
			return err
		}

		// update the recurring transaction
		if _, err := UpdateRecurring(ctx, &recurringTransaction); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":                err,
				"recurringTransaction": recurringTransaction,
			}).Error("error updating recurring transaction after adding the transaction. a duplicate transaction may be created in the future.")
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":                 err,
			"recurringTransactions": recurringTransactions,
		}).Error("error committing all recurring transactions and generated transactions")
		return err
	}

	return nil
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

	if err := validateRecurringTransaction(*transaction); err != nil {
		return nil, err
	}

	valid, err := userOwnsAccount(c, transaction.Transaction.AccountId)
	if err != nil || !valid {
		return nil, constants.Forbidden
	}

	tdb := recurringToDB(*transaction)
	var id int
	err = db.QueryRow("INSERT INTO recurringTransactions(name, nextOccurs, category, amount, note, accountId, scheduleType, secondsBetween, dayOf, secondsBeforeToPost) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id", tdb.Name, tdb.NextOccurs, tdb.Category, tdb.Amount, tdb.Note, tdb.AccountId, tdb.ScheduleType, tdb.SecondsBetween, tdb.DayOf, tdb.SecondsBeforeToPost).Scan(&id)
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

func UpdateRecurring(c context.Context, transaction *RecurringTransaction) (*RecurringTransaction, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	if err := validateRecurringTransaction(*transaction); err != nil {
		return nil, err
	}

	tdb := recurringToDB(*transaction)
	_, err = db.Exec("UPDATE recurringTransactions SET name = $1, nextOccurs = $2, category = $3, amount = $4, note = $5, accountId = $6, scheduleType = $7, secondsBetween = $8, dayOf = $9, secondsBeforeToPost = $10", tdb.Name, tdb.NextOccurs, tdb.Category, tdb.Amount, tdb.Note, tdb.AccountId, tdb.ScheduleType, tdb.SecondsBetween, tdb.DayOf, tdb.SecondsBeforeToPost)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":                  err,
			"recurringTransactionDB": tdb,
			"recurringTransaction":   transaction,
		}).Errorf("failed to update recurring transaction row")
		return nil, err
	}

	return transaction, nil
}

func validateRecurringTransaction(tr RecurringTransaction) error {
	switch tr.ScheduleType {
	case constants.FIXED_INTERVAL:
		if tr.SecondsBetween != nil && *tr.SecondsBetween != 0 {
			return nil
		}
	case constants.FIXED_DAY_WEEK, constants.FIXED_DAY_MONTH, constants.FIXED_DAY_YEAR:
		if tr.DayOf != nil {
			return nil
		}
	}

	return constants.BadRequest
}

func getNextRun(tr RecurringTransaction) (time.Time, error) {
	switch tr.ScheduleType {
	case constants.FIXED_INTERVAL:
		return tr.Transaction.Date.Add(time.Duration(*tr.SecondsBetween) * time.Second), nil
	case constants.FIXED_DAY_WEEK:
		currDay := util.WeekdayToInt(tr.Transaction.Date.Weekday())
		daysToAdd := *tr.DayOf - currDay
		if currDay == daysToAdd {
			daysToAdd += 7
		}
		return tr.Transaction.Date.Add(time.Hour * time.Duration(24*daysToAdd)), nil
	case constants.FIXED_DAY_MONTH:
		desiredDay := *tr.DayOf

		// at earliest, the next transaction must be after minDay in minMonth and minYear.
		// usually, it would be in the next month, unless the minDay is less than
		// min(desired day, days in minMonth), where the next transaction can be in minMonth.
		minDay := tr.Transaction.Date.Day()
		minMonth := tr.Transaction.Date.Month()
		minYear := tr.Transaction.Date.Year()
		daysInMinMonth := util.DaysIn(minMonth, minYear)

		// min of (desired day, days in minMonth) is the day the transaction should be scheduled
		// in this month. if that is greater than minDay/minMonth/minYear, just schedule
		// the next transaction in minMonth
		correctDayForMinMonth := util.Min(desiredDay, daysInMinMonth)
		if correctDayForMinMonth > minDay {
			return tr.Transaction.Date.Add(time.Hour * time.Duration(24*(correctDayForMinMonth-minDay))), nil
		}

		// minMonth didn't work out, so schedule for the next month
		year := minYear
		if minMonth == time.December {
			year += 1
		}
		month := minMonth + 1
		daysInMonth := util.DaysIn(month, year)
		correctDay := util.Min(desiredDay, daysInMonth)
		return time.Date(year, month, correctDay, 0, 0, 0, 0, time.UTC), nil
	case constants.FIXED_DAY_YEAR:
		// just need to check if it can run this year or must be next year
		desiredDay := *tr.DayOf
		minDayOfYear := tr.Transaction.Date.YearDay()
		if desiredDay > minDayOfYear {
			return tr.Transaction.Date.Add(time.Hour * time.Duration(24*(desiredDay-minDayOfYear))), nil
		}

		// increment the year
		newDate := time.Date(tr.Transaction.Date.Year()+1, tr.Transaction.Date.Month(), tr.Transaction.Date.Day(), 0, 0, 0, 0, time.UTC)

		// it is possible that newDate has the correct YearDay. consider starting with March 1 on a leap year. suppose it is
		// day x of that year. add 1 to the year. there are now less days before March 1, since Feb only has 28 days.
		if newDate.YearDay() == desiredDay {
			return newDate, nil
		}

		return newDate.Add(time.Hour * time.Duration(24*(desiredDay-newDate.YearDay()))), nil
	}

	err := fmt.Errorf("Unknown schedule type for recurring transaction: %s", tr.ScheduleType)
	logrus.WithFields(logrus.Fields{
		"error":                err,
		"recurringTransaction": tr,
	}).Error("unknown schedule type for recurring transaction")
	return time.Time{}, err
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
		Id: transaction.Id,
		Transaction: Transaction{
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
