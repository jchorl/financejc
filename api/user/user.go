package user

import (
	"context"
	"database/sql"

	"github.com/Sirupsen/logrus"
	"github.com/lib/pq"

	"github.com/jchorl/financejc/api/util"
	"github.com/jchorl/financejc/constants"
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
	err = db.QueryRow("SELECT email, google_id FROM users WHERE id = $1", userID).Scan(&email, &googleID)
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

// GetAll queries for all users
func GetAll(c context.Context) ([]User, error) {
	if !util.IsAdminRequest(c) {
		return nil, constants.ErrForbidden
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	users := []User{}
	rows, err := db.Query("SELECT id, google_id, email FROM users")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to fetch all users")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user userDB
		if err := rows.Scan(&user.ID, &user.GoogleID, &user.Email); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("failed to scan into user")
			return nil, err
		}

		users = append(users, fromDB(user))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to get users from rows")
		return nil, err
	}

	return users, nil
}

// BatchImport batch imports users
func BatchImport(c context.Context, users []User) error {
	if !util.IsAdminRequest(c) {
		return constants.ErrForbidden
	}

	db, err := util.SQLDBFromContext(c)
	if err != nil {
		return err
	}

	txn, err := db.Begin()
	if err != nil {
		logrus.WithError(err).Error("unable to begin transaction when batch inserting users")
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn("users", "id", "google_id", "email"))
	if err != nil {
		logrus.WithError(err).Error("unable to begin copy in when batch inserting users")
		return err
	}

	for _, user := range users {
		// skip admin user because they are already signed up
		if user.Email == constants.AdminEmail {
			continue
		}
		udb := toDB(user)
		_, err = stmt.Exec(udb.ID, udb.GoogleID, udb.Email)
		if err != nil {
			logrus.WithError(err).Error("unable to exec user copy when batch inserting users")
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		logrus.WithError(err).Error("unable to exec batch user copy when batch inserting users")
		return err
	}

	err = stmt.Close()
	if err != nil {
		logrus.WithError(err).Error("unable to close user copy when batch inserting users")
		return err
	}

	err = txn.Commit()
	if err != nil {
		logrus.WithError(err).Error("unable to commit user copy when batch inserting users")
		return err
	}

	return nil
}

// FindOrCreateByGoogleID finds a user with the given googleId, otherwise it creates one and returns it
func FindOrCreateByGoogleID(c context.Context, googleID, email string) (User, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return User{}, err
	}

	var id uint
	err = db.QueryRow("SELECT id FROM users WHERE google_id = $1", googleID).Scan(&id)
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
	err = db.QueryRow("INSERT INTO users (google_id, email) VALUES($1, $2) RETURNING id", udb.GoogleID, udb.Email).Scan(&id)
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

func fromDB(user userDB) User {
	return User{
		ID:       user.ID,
		Email:    user.Email,
		GoogleID: util.FromNullStringNonEmpty(user.GoogleID),
	}
}
