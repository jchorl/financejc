package handlers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"github.com/jchorl/financejc/api/account"
)

// GetAccounts fetches accounts
func GetAccounts(c echo.Context) error {
	accounts, err := account.Get(toContext(c))
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, accounts)
}

// NewAccount creates a new account
func NewAccount(c echo.Context) error {
	acc := new(account.Account)
	if err := c.Bind(acc); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("error parsing request to create account")
		return writeError(c, err)
	}

	acc, err := account.New(toContext(c), acc)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, acc)
}

// UpdateAccount updates an account
func UpdateAccount(c echo.Context) error {
	acc := new(account.Account)
	if err := c.Bind(acc); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("error parsing request to update account")
		return writeError(c, err)
	}

	acc, err := account.Update(toContext(c), acc)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, acc)
}

// DeleteAccount deletes an account
func DeleteAccount(c echo.Context) error {
	accountID, err := idFromParam(c, "accountId")
	if err != nil {
		return writeError(c, err)
	}

	if err := account.Delete(toContext(c), accountID); err != nil {
		return writeError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
