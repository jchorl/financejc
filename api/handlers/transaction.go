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
		return writeError(c, constants.ErrBadRequest)
	}

	next := c.QueryParam("start")
	transactions, err := transaction.Get(toContext(c), accountId, next)
	if err != nil {
		return writeError(c, err)
	}

	return writePaginatedEntity(c, transactions)
}

func GetRecurringTransactions(c echo.Context) error {
	accountIdStr := c.Param("accountId")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"context":   c,
			"accountId": accountIdStr,
		})
		return writeError(c, constants.ErrBadRequest)
	}

	transactions, err := transaction.GetRecurring(toContext(c), accountId)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, transactions)
}

func GetTemplates(c echo.Context) error {
	accountIdStr := c.Param("accountId")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"context":   c,
			"accountId": accountIdStr,
		})
		return writeError(c, constants.ErrBadRequest)
	}

	transactions, err := transaction.GetTemplates(toContext(c), accountId)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, transactions)
}

func NewTransaction(c echo.Context) error {
	tr := new(transaction.Transaction)
	if err := c.Bind(tr); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("unable to parse request to create transaction")
		return writeError(c, constants.ErrBadRequest)
	}

	accountIdStr := c.Param("accountId")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"context":   c,
			"accountId": accountIdStr,
		})
		return writeError(c, constants.ErrBadRequest)
	}

	tr.AccountID = accountId
	tr, err = transaction.New(toContext(c), tr)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, tr)
}

func NewRecurringTransaction(c echo.Context) error {
	tr := new(transaction.RecurringTransaction)
	if err := c.Bind(tr); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("unable to parse request to create recurring transaction")
		return writeError(c, constants.ErrBadRequest)
	}

	accountIdStr := c.Param("accountId")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"context":   c,
			"accountId": accountIdStr,
		})
		return writeError(c, constants.ErrBadRequest)
	}

	tr.Transaction.AccountID = accountId
	tr, err = transaction.NewRecurring(toContext(c), tr)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, tr)
}

func NewTemplate(c echo.Context) error {
	tr := new(transaction.Template)
	if err := c.Bind(tr); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("unable to parse request to create transaction template")
		return writeError(c, constants.ErrBadRequest)
	}

	accountIdStr := c.Param("accountId")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"context":   c,
			"accountId": accountIdStr,
		})
		return writeError(c, constants.ErrBadRequest)
	}

	tr.AccountID = accountId
	tr, err = transaction.NewTemplate(toContext(c), tr)
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
		return writeError(c, constants.ErrBadRequest)
	}

	tr, err := transaction.Update(toContext(c), tr)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, tr)
}

func UpdateRecurringTransaction(c echo.Context) error {
	tr := new(transaction.RecurringTransaction)
	if err := c.Bind(tr); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("unable to parse request to update recurring transaction")
		return writeError(c, constants.ErrBadRequest)
	}

	tr, err := transaction.UpdateRecurring(toContext(c), tr)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, tr)
}

func UpdateTemplate(c echo.Context) error {
	tr := new(transaction.Template)
	if err := c.Bind(tr); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("unable to parse request to update transaction template")
		return writeError(c, constants.ErrBadRequest)
	}

	tr, err := transaction.UpdateTemplate(toContext(c), tr)
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
		return writeError(c, constants.ErrBadRequest)
	}

	err = transaction.Delete(toContext(c), transactionId)
	if err != nil {
		return writeError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func DeleteRecurringTransaction(c echo.Context) error {
	transactionIdStr := c.Param("recurringTransactionId")
	transactionId, err := strconv.Atoi(transactionIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":                  err,
			"recurringTransactionId": transactionIdStr,
		}).Error("error parsing recurring transaction ID to int")
		return writeError(c, constants.ErrBadRequest)
	}

	err = transaction.DeleteRecurring(toContext(c), transactionId)
	if err != nil {
		return writeError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func DeleteTemplate(c echo.Context) error {
	transactionIdStr := c.Param("templateId")
	transactionId, err := strconv.Atoi(transactionIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"templateId": transactionIdStr,
		}).Error("error parsing transaction template ID to int")
		return writeError(c, constants.ErrBadRequest)
	}

	err = transaction.DeleteTemplate(toContext(c), transactionId)
	if err != nil {
		return writeError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func QueryES(ctx echo.Context) error {
	accountIdStr := ctx.Param("accountId")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"context":   ctx,
			"accountId": accountIdStr,
		})
		return writeError(ctx, constants.ErrBadRequest)
	}

	query := transaction.Query{
		Field:     ctx.QueryParam("field"),
		Value:     ctx.QueryParam("value"),
		AccountID: accountId,
	}
	transactions, err := transaction.QueryES(toContext(ctx), query)
	if err != nil {
		return writeError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, transactions)
}

func PushAllToES(ctx echo.Context) error {
	if err := transaction.PushAllToES(toContext(ctx)); err != nil {
		return writeError(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}
