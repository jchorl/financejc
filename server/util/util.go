package util

import (
	"database/sql"
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
