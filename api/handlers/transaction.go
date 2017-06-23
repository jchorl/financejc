package handlers

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"github.com/jchorl/financejc/api/transaction"
	"github.com/jchorl/financejc/constants"
)

// GetTransactions fetches transactions
func GetTransactions(c echo.Context) error {
	accountID, err := idFromParam(c, "accountId")
	if err != nil {
		return writeError(c, err)
	}

	next := c.QueryParam("start")
	transactions, err := transaction.Get(toContext(c), accountID, next)
	if err != nil {
		return writeError(c, err)
	}

	return writePaginatedEntity(c, transactions)
}

// GetSummary fetches all transactions since a given timestamp
func GetSummary(c echo.Context) error {
	sinceStr := c.QueryParam("since")
	since, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"context":  c,
			"sinceStr": sinceStr,
		})
		writeError(c, constants.ErrBadRequest)
		return constants.ErrBadRequest
	}

	transactions, err := transaction.Summary(toContext(c), since)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, transactions)
}

// GetRecurringTransactions fetches recurring transactions
func GetRecurringTransactions(c echo.Context) error {
	accountID, err := idFromParam(c, "accountId")
	if err != nil {
		return writeError(c, err)
	}

	transactions, err := transaction.GetRecurring(toContext(c), accountID)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, transactions)
}

// GetTemplates fetches templates for transactions
func GetTemplates(c echo.Context) error {
	accountID, err := idFromParam(c, "accountId")
	if err != nil {
		return writeError(c, err)
	}

	transactions, err := transaction.GetTemplates(toContext(c), accountID)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, transactions)
}

// NewTransaction creates a new transaction
func NewTransaction(c echo.Context) error {
	tr := new(transaction.Transaction)
	if err := c.Bind(tr); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("unable to parse request to create transaction")
		return writeError(c, constants.ErrBadRequest)
	}

	accountID, err := idFromParam(c, "accountId")
	if err != nil {
		return writeError(c, err)
	}

	tr.AccountID = accountID
	tr, err = transaction.New(toContext(c), tr)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, tr)
}

// NewRecurringTransaction creates a new recurring transaction
func NewRecurringTransaction(c echo.Context) error {
	tr := new(transaction.RecurringTransaction)
	if err := c.Bind(tr); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("unable to parse request to create recurring transaction")
		return writeError(c, constants.ErrBadRequest)
	}

	accountID, err := idFromParam(c, "accountId")
	if err != nil {
		return writeError(c, err)
	}

	tr.Transaction.AccountID = accountID
	tr, err = transaction.NewRecurring(toContext(c), tr)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, tr)
}

// NewTemplate creates a new template
func NewTemplate(c echo.Context) error {
	tr := new(transaction.Template)
	if err := c.Bind(tr); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("unable to parse request to create transaction template")
		return writeError(c, constants.ErrBadRequest)
	}

	accountID, err := idFromParam(c, "accountId")
	if err != nil {
		return writeError(c, err)
	}

	tr.AccountID = accountID
	tr, err = transaction.NewTemplate(toContext(c), tr)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, tr)
}

// UpdateTransaction updates a transaction
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

// UpdateRecurringTransaction updates a recurring transaction
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

// UpdateTemplate updates a template
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

// DeleteTransaction deletes a transaction
func DeleteTransaction(c echo.Context) error {
	transactionID, err := idFromParam(c, "transactionId")
	if err != nil {
		return writeError(c, err)
	}

	err = transaction.Delete(toContext(c), transactionID)
	if err != nil {
		return writeError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteRecurringTransaction deletes a recurring transaction
func DeleteRecurringTransaction(c echo.Context) error {
	recurringTransactionID, err := idFromParam(c, "recurringTransactionId")

	err = transaction.DeleteRecurring(toContext(c), recurringTransactionID)
	if err != nil {
		return writeError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteTemplate deletes a template
func DeleteTemplate(c echo.Context) error {
	templateID, err := idFromParam(c, "templateId")

	err = transaction.DeleteTemplate(toContext(c), templateID)
	if err != nil {
		return writeError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// QueryES parses a query and queries elasticsearch for matching transactions
func QueryES(c echo.Context) error {
	accountID, err := idFromParam(c, "accountId")
	if err != nil {
		return writeError(c, err)
	}

	query := transaction.Query{
		Field:     c.QueryParam("field"),
		Value:     c.QueryParam("value"),
		AccountID: accountID,
	}
	transactions, err := transaction.QueryES(toContext(c), query)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, transactions)
}

// Search does a general search against all transactions
func Search(c echo.Context) error {
	transactions, err := transaction.SearchES(toContext(c), c.QueryParam("value"))
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, transactions)
}

// PushAllToES destroys the elasticsearch index and repushes all transactions
func PushAllToES(ctx echo.Context) error {
	if err := transaction.PushAllToES(toContext(ctx)); err != nil {
		return writeError(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// GenRecurringTransactions generates all recurring transactions
func GenRecurringTransactions(ctx echo.Context) error {
	if err := transaction.GenRecurringTransactions(toContext(ctx)); err != nil {
		return writeError(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}
