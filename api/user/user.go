package user

import (
	"context"
	"database/sql"

	"github.com/Sirupsen/logrus"

	"github.com/jchorl/financejc/api/util"
)

// User represents a user
type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	GoogleID string `json:"-"`
}

type userDB struct {
	ID       uint
	Email    string
	GoogleID sql.NullString
}

// Get gets a user from the ID baked into the context
func Get(c context.Context) (User, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return User{}, err
	}

	userID, err := util.UserIDFromContext(c)
	if err != nil {
		return User{}, err
	}

	var email, googleID string
	err = db.QueryRow("SELECT email, googleId FROM users WHERE id = $1", userID).Scan(&email, &googleID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"googleId": googleID,
		}).Error("failed to select user from users table")
		return User{}, err
	}

	return User{
		ID:       userID,
		Email:    email,
		GoogleID: googleID,
	}, nil
}

// FindOrCreateByGoogleID finds a user with the given googleId, otherwise it creates one and returns it
func FindOrCreateByGoogleID(c context.Context, googleID, email string) (User, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return User{}, err
	}

	var id uint
	err = db.QueryRow("SELECT id FROM users WHERE googleId = $1", googleID).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"googleId": googleID,
		}).Error("failed to select id from users table")
		return User{}, err
	} else if err == nil {
		return User{
			ID:       id,
			Email:    email,
			GoogleID: googleID,
		}, nil
	}

	user := User{
		Email:    email,
		GoogleID: googleID,
	}
	udb := toDB(user)
	err = db.QueryRow("INSERT INTO users (googleId, email) VALUES($1, $2) RETURNING id", udb.GoogleID, udb.Email).Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"googleId": googleID,
			"user":     user,
			"userDb":   udb,
		}).Error("failed to insert user")
		return User{}, err
	}

	user.ID = id
	return user, nil
}

func toDB(user User) *userDB {
	return &userDB{
		ID:       user.ID,
		Email:    user.Email,
		GoogleID: util.ToNullStringNonEmpty(user.GoogleID),
	}
}

func fromDB(user userDB) *User {
	return &User{
		ID:       user.ID,
		Email:    user.Email,
		GoogleID: util.FromNullStringNonEmpty(user.GoogleID),
	}
}
