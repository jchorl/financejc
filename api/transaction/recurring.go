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

// RecurringTransaction is a template for a transaction that gets automatically generated
type RecurringTransaction struct {
	ID                  int         `json:"id,omitempty"`
	Transaction         Transaction `json:"transaction"`
	ScheduleType        string      `json:"scheduleType"`
	SecondsBetween      *int        `json:"secondsBetween"`
	DayOf               *int        `json:"dayOf"`
	SecondsBeforeToPost int         `json:"secondsBeforeToPost"`
}

type recurringTransactionDB struct {
	ID         int
	Name       string
	NextOccurs time.Time
	Category   string
	Amount     int
	Note       string
	AccountID  int

	ScheduleType        string
	SecondsBetween      sql.NullInt64
	DayOf               sql.NullInt64
	SecondsBeforeToPost int
}

// GenRecurringTransactions generates transactions from recurring transactions
func GenRecurringTransactions(ctx context.Context) error {
	logrus.Debug("running GenRecurringTransactions")
	db, err := util.SQLDBFromContext(ctx)
	if err != nil {
		return err
	}

	recurringTransactions, err := getRecurringToPost(db)
	if err != nil {
		return err
	}

	logrus.Debugf("got %d recurring transactions to post", len(recurringTransactions))

	tx, err := db.Begin()
	if err != nil {
		logrus.WithError(err).Error("could not begin transaction to generate recurring transactions")
		return err
	}

	// replace the sql Db in the context with the sql Tx
	ctx = context.WithValue(ctx, constants.CTX_DB, tx)

	for _, recurringTransaction := range recurringTransactions {
		logrus.WithField("recurringTransaction", recurringTransaction).Debug("about to generate recurring transaction")
		if err := generateFromRecurringAndUpdateRecurring(ctx, recurringTransaction); err != nil {
			tx.Rollback()
			return err
		}
		logrus.Debug("generated")
	}

	if err := tx.Commit(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":                 err,
			"recurringTransactions": recurringTransactions,
		}).Error("error committing all recurring transactions and generated transactions")
		tx.Rollback()
		return err
	}

	logrus.Debug("committed and returning")
	return nil
}

// GetRecurring fetches the recurring transactions for a given account
func GetRecurring(c context.Context, accountID int) ([]RecurringTransaction, error) {
	transactions := []RecurringTransaction{}
	db, err := util.DBFromContext(c)
	if err != nil {
		return transactions, err
	}

	valid, err := userOwnsAccount(c, accountID)
	if err != nil || !valid {
		return transactions, constants.Forbidden
	}

	rows, err := db.Query("SELECT id, name, nextOccurs, category, amount, note, accountId, scheduleType, secondsBetween, dayOf, secondsBeforeToPost FROM recurringTransactions WHERE accountId = $1", accountID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountID,
		}).Error("failed to fetch recurring transactions")
		return transactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction recurringTransactionDB
		if err := rows.Scan(&transaction.ID, &transaction.Name, &transaction.NextOccurs, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.AccountID, &transaction.ScheduleType, &transaction.SecondsBetween, &transaction.DayOf, &transaction.SecondsBeforeToPost); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"accountId": accountID,
			}).Error("failed to scan into recurring transaction")
			return transactions, err
		}

		transactions = append(transactions, recurringFromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountID,
		}).Error("failed to get recurring transactions from rows")
		return transactions, err
	}

	return transactions, nil
}

// NewRecurring creates a new recurring transaction
func NewRecurring(c context.Context, transaction *RecurringTransaction) (*RecurringTransaction, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	if err := validateRecurringTransaction(*transaction); err != nil {
		return nil, err
	}

	valid, err := userOwnsAccount(c, transaction.Transaction.AccountID)
	if err != nil || !valid {
		return nil, constants.Forbidden
	}

	transaction.Transaction.Date, err = getNextRun(transaction, true)
	if err != nil {
		return nil, err
	}

	tdb := recurringToDB(*transaction)
	var id int
	err = db.QueryRow("INSERT INTO recurringTransactions(name, nextOccurs, category, amount, note, accountId, scheduleType, secondsBetween, dayOf, secondsBeforeToPost) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id", tdb.Name, tdb.NextOccurs, tdb.Category, tdb.Amount, tdb.Note, tdb.AccountID, tdb.ScheduleType, tdb.SecondsBetween, tdb.DayOf, tdb.SecondsBeforeToPost).Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":                  err,
			"recurringTransactionDB": tdb,
			"recurringTransaction":   transaction,
		}).Errorf("failed to insert recurring transaction row")
		return nil, err
	}

	transaction.ID = id
	return transaction, nil
}

// UpdateRecurring updates a recurring transaction
func UpdateRecurring(c context.Context, transaction *RecurringTransaction) (*RecurringTransaction, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	transaction.Transaction.Date, err = getNextRun(transaction, true)
	if err != nil {
		return nil, err
	}

	if err := validateRecurringTransaction(*transaction); err != nil {
		return nil, err
	}

	tdb := recurringToDB(*transaction)
	_, err = db.Exec("UPDATE recurringTransactions SET name = $1, nextOccurs = $2, category = $3, amount = $4, note = $5, accountId = $6, scheduleType = $7, secondsBetween = $8, dayOf = $9, secondsBeforeToPost = $10", tdb.Name, tdb.NextOccurs, tdb.Category, tdb.Amount, tdb.Note, tdb.AccountID, tdb.ScheduleType, tdb.SecondsBetween, tdb.DayOf, tdb.SecondsBeforeToPost)
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

// DeleteRecurring deletes a recurring transaction
func DeleteRecurring(ctx context.Context, transactionID int) error {
	db, err := util.DBFromContext(ctx)
	if err != nil {
		return err
	}

	valid, err := userOwnsRecurringTransaction(ctx, transactionID)
	if err != nil || !valid {
		return constants.Forbidden
	}

	_, err = db.Exec("DELETE FROM recurringTransactions WHERE id = $1", transactionID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":                  err,
			"recurringTransactionID": transactionID,
		}).Errorf("could not delete recurring transaction")
		return err
	}

	return nil
}

func userOwnsRecurringTransaction(c context.Context, recurringTransaction int) (bool, error) {
	userID, err := util.UserIdFromContext(c)
	if err != nil {
		return false, err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return false, err
	}

	var owner uint
	err = db.QueryRow("SELECT a.userId FROM accounts a JOIN recurringTransactions t ON t.accountId = a.id WHERE t.id = $1", recurringTransaction).Scan(&owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":                err,
			"userId":               userID,
			"recurringTransaction": recurringTransaction,
		}).Error("error checking owner of recurring transaction")
		return false, err
	}

	return owner == userID, nil
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

func getNextRun(tr *RecurringTransaction, allowSameDay bool) (time.Time, error) {
	switch tr.ScheduleType {
	case constants.FIXED_INTERVAL:
		if allowSameDay {
			return tr.Transaction.Date, nil
		}
		return tr.Transaction.Date.Add(time.Duration(*tr.SecondsBetween) * time.Second), nil
	case constants.FIXED_DAY_WEEK:
		currDay := util.WeekdayToInt(tr.Transaction.Date.Weekday())
		// force positive
		daysToAdd := (*tr.DayOf - currDay + 7) % 7
		if daysToAdd == 0 && !allowSameDay {
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
		if correctDayForMinMonth > minDay || correctDayForMinMonth == minDay && allowSameDay {
			return tr.Transaction.Date.Add(time.Hour * time.Duration(24*(correctDayForMinMonth-minDay))), nil
		}

		// minMonth didn't work out, so schedule for the next month
		year := minYear
		if minMonth == time.December {
			year++
		}
		month := minMonth + 1
		daysInMonth := util.DaysIn(month, year)
		correctDay := util.Min(desiredDay, daysInMonth)
		return time.Date(year, month, correctDay, 0, 0, 0, 0, time.UTC), nil
	case constants.FIXED_DAY_YEAR:
		// just need to check if it can run this year or must be next year
		desiredDay := *tr.DayOf
		minDayOfYear := tr.Transaction.Date.YearDay()
		if desiredDay > minDayOfYear || desiredDay == minDayOfYear && allowSameDay {
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

func getRecurringToPost(db util.DB) ([]RecurringTransaction, error) {
	// query for all recurring transactions where the next occurrance is within the time period before to post it to the account
	rows, err := db.Query("SELECT id, name, nextOccurs, category, amount, note, accountId, scheduleType, secondsBetween, dayOf, secondsBeforeToPost FROM recurringTransactions WHERE nextOccurs - interval '1 second' * secondsBeforeToPost <= NOW()")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to fetch recurring transactions to post to accounts")
		return nil, err
	}
	defer rows.Close()

	recurringTransactions := []RecurringTransaction{}
	for rows.Next() {
		var transaction recurringTransactionDB
		if err := rows.Scan(&transaction.ID, &transaction.Name, &transaction.NextOccurs, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.AccountID, &transaction.ScheduleType, &transaction.SecondsBetween, &transaction.DayOf, &transaction.SecondsBeforeToPost); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("failed to scan into recurring transaction to generate transaction")
			return nil, err
		}

		recurringTransactions = append(recurringTransactions, recurringFromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to get recurring transactions from rows to generate transactions")
		return nil, err
	}

	return recurringTransactions, nil
}

func generateFromRecurringAndUpdateRecurring(ctx context.Context, recurringTransaction RecurringTransaction) error {
	// keep generating until it is too early to post the next transaction
	now := time.Now()
	for recurringTransaction.Transaction.Date.Add(time.Second * time.Duration(-recurringTransaction.SecondsBeforeToPost)).Before(now) {
		if _, err := newWithoutVerifyingAccountOwnership(ctx, &recurringTransaction.Transaction); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":                err,
				"recurringTransaction": recurringTransaction,
			}).Error("error adding a new transaction for a recurring transaction")
			return err
		}

		// calculate when the transaction should next run
		var err error
		recurringTransaction.Transaction.Date, err = getNextRun(&recurringTransaction, false)
		if err != nil {
			return err
		}
	}

	// update the recurring transaction
	if _, err := UpdateRecurring(ctx, &recurringTransaction); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":                err,
			"recurringTransaction": recurringTransaction,
		}).Error("error updating recurring transaction after adding the transaction. a duplicate transaction may be created in the future.")
		return err
	}

	return nil
}

func recurringToDB(transaction RecurringTransaction) *recurringTransactionDB {
	return &recurringTransactionDB{
		ID:         transaction.ID,
		Name:       transaction.Transaction.Name,
		NextOccurs: transaction.Transaction.Date,
		Category:   transaction.Transaction.Category,
		Amount:     transaction.Transaction.Amount,
		Note:       transaction.Transaction.Note,
		AccountID:  transaction.Transaction.AccountID,

		ScheduleType:        transaction.ScheduleType,
		SecondsBetween:      util.ToNullInt(transaction.SecondsBetween),
		DayOf:               util.ToNullInt(transaction.DayOf),
		SecondsBeforeToPost: transaction.SecondsBeforeToPost,
	}
}

func recurringFromDB(transaction recurringTransactionDB) RecurringTransaction {
	return RecurringTransaction{
		ID: transaction.ID,
		Transaction: Transaction{
			Name:      transaction.Name,
			Date:      transaction.NextOccurs,
			Category:  transaction.Category,
			Amount:    transaction.Amount,
			Note:      transaction.Note,
			AccountID: transaction.AccountID,
		},

		ScheduleType:        transaction.ScheduleType,
		SecondsBetween:      util.FromNullInt(transaction.SecondsBetween),
		DayOf:               util.FromNullInt(transaction.DayOf),
		SecondsBeforeToPost: transaction.SecondsBeforeToPost,
	}
}
