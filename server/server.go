package server

import (
	"context"
	"database/sql"

	"github.com/jchorl/financejc/constants"

	_ "github.com/lib/pq"
)

type server struct {
	db *sql.DB
}

func NewServer(driver, address string) (server, error) {
	db, err := sql.Open(driver, address)
	if err != nil {
		return server{}, err
	}

	s := server{db}
	return s, nil
}

func (s server) DB() *sql.DB {
	return s.db
}

func (s server) Context() context.Context {
	c := context.Background()
	return context.WithValue(c, constants.CTX_DB, s.DB())
}

func (s server) ContextWithUser(user int) context.Context {
	return context.WithValue(s.Context(), constants.CTX_USER, user)
}
