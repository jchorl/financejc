package handlers

import (
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"github.com/jchorl/financejc/api/transaction"
	"github.com/jchorl/financejc/constants"
)

func GetTransactions(c echo.Context) error {
	accountIdStr := c.Param("accountId")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"context":   c,
			"accountId": accountIdStr,
		})
		return writeError(c, constants.BadRequest)
	}

	next := c.QueryParam("start")
	transactions, err := transaction.Get(c, accountId, next)
	if err != nil {
		return writeError(c, err)
	}

	return writePaginatedEntity(c, transactions)
}

func NewTransaction(c echo.Context) error {
	tr := new(transaction.Transaction)
	if err := c.Bind(tr); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("unable to parse request to create transaction")
		return writeError(c, constants.BadRequest)
	}

	accountIdStr := c.Param("accountId")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"context":   c,
			"accountId": accountIdStr,
		})
		return writeError(c, constants.BadRequest)
	}

	tr.Account = accountId
	tr, err = transaction.New(c, tr)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, tr)
}

func UpdateTransaction(c echo.Context) error {
	tr := new(transaction.Transaction)
	if err := c.Bind(tr); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("unable to parse request to update transaction")
		return writeError(c, constants.BadRequest)
	}

	tr, err := transaction.Update(c, tr)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, tr)
}

func DeleteTransaction(c echo.Context) error {
	transactionIdStr := c.Param("transactionId")
	transactionId, err := strconv.Atoi(transactionIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":         err,
			"transactionId": transactionIdStr,
		}).Error("error parsing transaction ID to int")
		return writeError(c, constants.BadRequest)
	}

	err = transaction.Delete(c, transactionId)
	if err != nil {
		return writeError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
