package util

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"

	"github.com/jchorl/financejc/constants"
)

// DB is an interface that both *sql.DB and *sql.Tx conform to
// Functions take the DB interface so they can be passed
// a *sql.DB or a *sql.Tx and use them the same way.
type DB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
}

var _ DB = (*sql.DB)(nil)
var _ DB = (*sql.Tx)(nil)

var daysBefore = [...]int32{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}

var weekdayToInt = map[time.Weekday]int{
	time.Sunday:    0,
	time.Monday:    1,
	time.Tuesday:   2,
	time.Wednesday: 3,
	time.Thursday:  4,
	time.Friday:    5,
	time.Saturday:  6,
}

// ToNullStringNonEmpty converts a string s to a sql.NullString, treating empty string as NULL
func ToNullStringNonEmpty(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// ToNullIntNonZero converts an int i to a sql.NullInt64, treating 0 as NULL
func ToNullIntNonZero(i int) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(i), Valid: i != 0}
}

// FromNullStringNonEmpty converts a sql.NullString s to a string, treating NULL as empty string
func FromNullStringNonEmpty(s sql.NullString) string {
	if !s.Valid {
		return ""
	}
	return s.String
}

// FromNullIntNonZero converts a sql.NullInt64 to an int, treating NULL as 0
func FromNullIntNonZero(i sql.NullInt64) int {
	if !i.Valid {
		return 0
	}
	return int(i.Int64)
}

// ToNullString turns a string pointer into a sql nullable string
func ToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

// FromNullString takes a sql nullable string and returns a string pointer
func FromNullString(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// ToNullInt turns an int pointer into a sql nullable int
func ToNullInt(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

// FromNullInt takes a sql nullable int and returns an int pointer
func FromNullInt(i sql.NullInt64) *int {
	if !i.Valid {
		return nil
	}
	conv := int(i.Int64)
	return &conv
}

// UserIDFromContext pulls a user ID from a context
func UserIDFromContext(c context.Context) (uint, error) {
	userID, ok := c.Value(constants.CtxUserID).(uint)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"context": c,
		}).Error("Unable to get userId from context")
		return 0, errors.New("Unable to get userId from context")
	}
	return userID, nil
}

// DBFromContext pulls a database connection from a context
func DBFromContext(c context.Context) (DB, error) {
	db := c.Value(constants.CtxDB).(DB)
	if db == nil {
		logrus.WithFields(logrus.Fields{
			"context": c,
		}).Error("Unable to get DB from context")
		return nil, errors.New("Unable to get DB from context")
	}

	return db, nil
}

// ESFromContext pulls an elasticsearch connection from a context
func ESFromContext(c context.Context) (*elastic.Client, error) {
	es := c.Value(constants.CtxES)
	if es == nil {
		logrus.WithFields(logrus.Fields{
			"context": c,
		}).Error("Unable to get ES client from context")
		return nil, errors.New("Unable to get ES client from context")
	}

	parsed, ok := es.(*elastic.Client)
	if !ok || parsed == nil {
		logrus.WithFields(logrus.Fields{
			"context": c,
		}).Error("Unable to get ES client from context, could not cast val to client")
		return nil, errors.New("Unable to get ES client from context")
	}

	return parsed, nil
}

// SQLDBFromContext returns a *sql.DB from a context
// make sure to call with a context that will have a *sql.db and not a DB
func SQLDBFromContext(c context.Context) (*sql.DB, error) {
	db, ok := c.Value(constants.CtxDB).(*sql.DB)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"context": c,
		}).Error("Unable to get SQL DB from context")
		return nil, errors.New("Unable to get SQL DB from context")
	}

	return db, nil
}

// Min returns the min of two ints
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// DaysIn returns the number of days in a month
func DaysIn(m time.Month, year int) int {
	if m == time.February && isLeap(year) {
		return 29
	}
	return int(daysBefore[m] - daysBefore[m-1])
}

// WeekdayToInt converts a weekday to an integer value
func WeekdayToInt(d time.Weekday) int {
	return weekdayToInt[d]
}

// UserOwnsAccount checks if a user owns an account
func UserOwnsAccount(c context.Context, accountID int) (bool, error) {
	userID, err := UserIDFromContext(c)
	if err != nil {
		return false, err
	}

	db, err := DBFromContext(c)
	if err != nil {
		return false, err
	}

	var owner uint
	err = db.QueryRow("SELECT userId FROM accounts WHERE id = $1", accountID).Scan(&owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"userId":  userID,
			"account": accountID,
		}).Error("error checking owner of account")
		return false, err
	}

	return owner == userID, nil
}

// UserOwnsTransaction checks if a user owns a transaction
func UserOwnsTransaction(c context.Context, transactionID int) (bool, error) {
	userID, err := UserIDFromContext(c)
	if err != nil {
		return false, err
	}

	db, err := DBFromContext(c)
	if err != nil {
		return false, err
	}

	var owner uint
	err = db.QueryRow("SELECT a.userId FROM accounts a JOIN transactions t ON t.accountId = a.id WHERE t.id = $1", transactionID).Scan(&owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":       err,
			"userId":      userID,
			"transaction": transactionID,
		}).Error("error checking owner of transaction")
		return false, err
	}

	return owner == userID, nil
}
