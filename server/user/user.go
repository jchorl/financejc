package user

import (
	"context"
	"database/sql"

	"github.com/jchorl/financejc/constants"
	"github.com/jchorl/financejc/server/util"

	"github.com/Sirupsen/logrus"
)

type User struct {
	Id       int
	GoogleId string
}

type userDB struct {
	Id       int
	GoogleId sql.NullString
}

func FindOrCreateByGoogleId(c context.Context, googleId string) (int, error) {
	db := c.Value(constants.CTX_DB).(util.DB)

	var id int
	err := db.QueryRow("SELECT id FROM users where googleId = $1", googleId).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		logrus.WithFields(logrus.Fields{
			"Error":     err,
			"Google ID": googleId,
		}).Error("failed to select id from users table")
		return -1, err
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
			"Error":     err,
			"Google ID": googleId,
			"User":      user,
			"UserDB":    udb,
		}).Error("failed to insert user")
		return -1, err
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
