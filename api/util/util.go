package util

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"

	"github.com/jchorl/financejc/constants"
)

// Both *sql.DB and *sql.Tx conform to DB
type DB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
}

var _ DB = (*sql.DB)(nil)
var _ DB = (*sql.Tx)(nil)

func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func ToNullInt(i int) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(i), Valid: i != 0}
}

func FromNullString(s sql.NullString) string {
	if !s.Valid {
		return ""
	}
	return s.String
}

func FromNullInt(i sql.NullInt64) int {
	if !i.Valid {
		return 0
	}
	return int(i.Int64)
}

func UserIdFromContext(c context.Context) (uint, error) {
	user := c.Value("user").(*jwt.Token)
	if user == nil {
		logrus.WithFields(logrus.Fields{
			"context": c,
		}).Error("Unable to get user from context")
		return 0, errors.New("Unable to get user from context")
	}

	claims := user.Claims.(jwt.MapClaims)
	userIdF := claims["sub"].(float64)
	userId := uint(userIdF)
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
