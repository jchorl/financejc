package user

import (
	"context"
	"database/sql"

	"github.com/Sirupsen/logrus"

	"github.com/jchorl/financejc/api/util"
)

type User struct {
	Id       uint   `json:"id"`
	Email    string `json:"email"`
	GoogleId string `json:"-"`
}

type userDB struct {
	Id       uint
	Email    string
	GoogleId sql.NullString
}

func Get(c context.Context) (User, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return User{}, err
	}

	userId, err := util.UserIdFromContext(c)
	if err != nil {
		return User{}, err
	}

	var email, googleId string
	err = db.QueryRow("SELECT email, googleId FROM users WHERE id = $1", userId).Scan(&email, &googleId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"googleId": googleId,
		}).Error("failed to select user from users table")
		return User{}, err
	}

	return User{
		Id:       userId,
		Email:    email,
		GoogleId: googleId,
	}, nil
}

func FindOrCreateByGoogleId(c context.Context, googleId, email string) (User, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return User{}, err
	}

	var id uint
	err = db.QueryRow("SELECT id FROM users WHERE googleId = $1", googleId).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"googleId": googleId,
		}).Error("failed to select id from users table")
		return User{}, err
	} else if err == nil {
		return User{
			Id:       id,
			Email:    email,
			GoogleId: googleId,
		}, nil
	}

	user := User{
		Email:    email,
		GoogleId: googleId,
	}
	udb := toDB(user)
	err = db.QueryRow("INSERT INTO users (googleId, email) VALUES($1, $2) RETURNING id", udb.GoogleId, udb.Email).Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"googleId": googleId,
			"user":     user,
			"userDb":   udb,
		}).Error("failed to insert user")
		return User{}, err
	}

	user.Id = id
	return user, nil
}

func toDB(user User) *userDB {
	return &userDB{
		Id:       user.Id,
		Email:    user.Email,
		GoogleId: util.ToNullStringNonEmpty(user.GoogleId),
	}
}

func fromDB(user userDB) *User {
	return &User{
		Id:       user.Id,
		Email:    user.Email,
		GoogleId: util.FromNullStringNonEmpty(user.GoogleId),
	}
}
