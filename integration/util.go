// +build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/jchorl/financejc/api/account"
	"github.com/jchorl/financejc/api/user"
	"github.com/jchorl/financejc/constants"
)

func FreshDb(t *testing.T) *sql.DB {
	db, err := sql.Open(constants.DB_DRIVER, constants.DB_ADDRESS)
	assert.NoError(t, err, "unable to connect to db")

	tx, err := db.Begin()
	assert.NoError(t, err, "unable to begin transaction")

	rows, err := tx.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE'")
	assert.NoError(t, err, "unable to query for table names")
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		assert.NoError(t, err, "could not scan table name into string")

		// tx.Exec will reject using $1 for the table name due to sql injection
		// obviously manually formatting a string is not ideal, but this is done in a test
		_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s", name))
		assert.NoError(t, err, "could not delete from table %s", name)
	}
	err = rows.Err()
	assert.NoError(t, err, "rows error")

	err = tx.Commit()
	assert.NoError(t, err, "unable to commit transaction clearing the db")

	return db
}

func ContextWithUserAndDB(uid uint, db *sql.DB) context.Context {
	return context.WithValue(context.WithValue(context.Background(), constants.CTX_USER_ID, uid), constants.CTX_DB, db)
}

func NewUser(t *testing.T, ctx context.Context) uint {
	googleId := rand.Int()
	uid, err := user.FindOrCreateByGoogleId(ctx, strconv.Itoa(googleId))
	assert.NoError(t, err, "unable to create user")

	return uid
}

func NewAccount(t *testing.T, ctx context.Context) *account.Account {
	name := rand.Int()
	acc := &account.Account{
		Name:     strconv.Itoa(name),
		Currency: "USD",
	}

	acc, err := account.New(ctx, acc)
	assert.NoError(t, err, "unable to create account")

	return acc
}
