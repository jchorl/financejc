package user

import (
	"context"
	"database/sql"

	"github.com/jchorl/financejc/constants"

	"github.com/Sirupsen/logrus"
)

type User struct {
	Id       int
	GoogleId string
}

func FindOrCreateByGoogleId(c context.Context, googleId string) (int, error) {
	db := c.Value(constants.CTX_DB).(sql.DB)

	var id int
	err := db.QueryRow("SELECT id FROM users where googleId = $1", googleId).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		logrus.WithField("Error", err).Error("failed to select id from users table")
		return -1, err
	} else if err == nil {
		return id, nil
	}

	err = db.QueryRow("INSERT INTO users(googleId) VALUES($1) RETURNING id", googleId).Scan(&id)
	if err != nil {
		logrus.WithField("Error", err).Error("failed to insert user")
		return -1, err
	}

	return id, nil
}
