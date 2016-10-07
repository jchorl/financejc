package user

import (
	"context"
	"database/sql"

	"github.com/Sirupsen/logrus"

	"github.com/jchorl/financejc/api/util"
)

type User struct {
	Id       uint
	GoogleId string
}

type userDB struct {
	Id       uint
	GoogleId sql.NullString
}

func FindOrCreateByGoogleId(c context.Context, googleId string) (uint, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return 0, err
	}

	var id uint
	err = db.QueryRow("SELECT id FROM users WHERE googleId = $1", googleId).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"googleId": googleId,
		}).Error("failed to select id from users table")
		return 0, err
	} else if err == nil {
		return id, nil
	}

	user := User{
		GoogleId: googleId,
	}
	udb := toDB(user)
	err = db.QueryRow("INSERT INTO users(googleId) VALUES($1) RETURNING id", udb.GoogleId).Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"googleId": googleId,
			"user":     user,
			"userDb":   udb,
		}).Error("failed to insert user")
		return 0, err
	}

	return id, nil
}

func toDB(user User) *userDB {
	return &userDB{
		Id:       user.Id,
		GoogleId: util.ToNullString(user.GoogleId),
	}
}

func fromDB(user userDB) *User {
	return &User{
		Id:       user.Id,
		GoogleId: util.FromNullString(user.GoogleId),
	}
}
