package util

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/jchorl/financejc/constants"
)

// Both *sql.DB and *sql.Tx conform to DB
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

func ToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func ToNullInt(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

func FromNullString(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

func FromNullInt(i sql.NullInt64) *int {
	if !i.Valid {
		return nil
	}
	conv := int(i.Int64)
	return &conv
}

func UserIdFromContext(c context.Context) (uint, error) {
	userId, ok := c.Value(constants.CTX_USER_ID).(uint)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"context": c,
		}).Error("Unable to get user from context")
		return 0, errors.New("Unable to get user from context")
	}
	return userId, nil
}

func DBFromContext(c context.Context) (DB, error) {
	db := c.Value(constants.CTX_DB).(DB)
	if db == nil {
		logrus.WithFields(logrus.Fields{
			"context": c,
		}).Error("Unable to get DB from context")
		return nil, errors.New("Unable to get DB from context")
	}

	return db, nil
}

func SQLDBFromContext(c context.Context) (*sql.DB, error) {
	db, ok := c.Value(constants.CTX_DB).(*sql.DB)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"context": c,
		}).Error("Unable to get SQL DB from context")
		return nil, errors.New("Unable to get SQL DB from context")
	}

	return db, nil
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func DaysIn(m time.Month, year int) int {
	if m == time.February && isLeap(year) {
		return 29
	}
	return int(daysBefore[m] - daysBefore[m-1])
}

func WeekdayToInt(d time.Weekday) int {
	return weekdayToInt[d]
}
