package handlers

import (
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"github.com/jchorl/financejc/api/account"
)

func GetAccounts(c echo.Context) error {
	accounts, err := account.Get(toContext(c))
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, accounts)
}

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

func DeleteAccount(c echo.Context) error {
	accountIdStr := c.Param("accountId")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountIdStr,
		}).Error("error parsing account ID to int")
		return writeError(c, err)
	}

	if err := account.Delete(toContext(c), accountId); err != nil {
		return writeError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
