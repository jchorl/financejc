package transaction

import (
	"context"

	"github.com/Sirupsen/logrus"

	"github.com/jchorl/financejc/api/util"
	"github.com/jchorl/financejc/constants"
)

// Template represents a template to create a new transaction
// The date part of the embedded transaction is ignored
type Template struct {
	Transaction
	TemplateName string `json:"templateName"`
}

type templateDB struct {
	ID           int
	TemplateName string
	Name         string
	Category     string
	Amount       int
	Note         string
	AccountID    int
}

// GetTemplates fetches the templates for account accountID
func GetTemplates(c context.Context, accountID int) ([]Template, error) {
	transactions := []Template{}
	db, err := util.DBFromContext(c)
	if err != nil {
		return transactions, err
	}

	valid, err := userOwnsAccount(c, accountID)
	if err != nil || !valid {
		return transactions, constants.Forbidden
	}

	rows, err := db.Query("SELECT id, templateName, name, category, amount, note, accountId FROM templates WHERE accountId = $1", accountID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountID": accountID,
		}).Error("failed to fetch transaction templates")
		return transactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction templateDB
		if err := rows.Scan(&transaction.ID, &transaction.TemplateName, &transaction.Name, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.AccountID); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"accountID": accountID,
			}).Error("failed to scan into transaction template")
			return transactions, err
		}

		transactions = append(transactions, templateFromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountID": accountID,
		}).Error("failed to get transaction templates from rows")
		return transactions, err
	}

	return transactions, nil
}

// NewTemplate creates a new template
func NewTemplate(c context.Context, transaction *Template) (*Template, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	valid, err := userOwnsAccount(c, transaction.AccountID)
	if err != nil || !valid {
		return nil, constants.Forbidden
	}

	tdb := templateToDB(*transaction)
	var id int
	err = db.QueryRow("INSERT INTO templates(templateName, name, category, amount, note, accountId) VALUES($1, $2, $3, $4, $5, $6) RETURNING id", tdb.TemplateName, tdb.Name, tdb.Category, tdb.Amount, tdb.Note, tdb.AccountID).Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"templateDB": tdb,
			"template":   transaction,
		}).Errorf("failed to insert transaction template row")
		return nil, err
	}

	transaction.ID = id
	return transaction, nil
}

// UpdateTemplate updates a template
func UpdateTemplate(c context.Context, transaction *Template) (*Template, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	tdb := templateToDB(*transaction)
	_, err = db.Exec("UPDATE templates SET templateName = $1, name = $2, category = $3, amount = $4, note = $5, accountId = $6 WHERE id = $7", tdb.TemplateName, tdb.Name, tdb.Category, tdb.Amount, tdb.Note, tdb.AccountID, tdb.ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"templateDB": tdb,
			"template":   transaction,
		}).Errorf("failed to update transaction template row")
		return nil, err
	}

	return transaction, nil
}

// DeleteTemplate deletes a template
func DeleteTemplate(ctx context.Context, transactionID int) error {
	db, err := util.DBFromContext(ctx)
	if err != nil {
		return err
	}

	valid, err := userOwnsTemplate(ctx, transactionID)
	if err != nil || !valid {
		return constants.Forbidden
	}

	_, err = db.Exec("DELETE FROM templates WHERE id = $1", transactionID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"templateId": transactionID,
		}).Errorf("could not delete transaction template")
		return err
	}

	return nil
}

func userOwnsTemplate(c context.Context, template int) (bool, error) {
	userID, err := util.UserIdFromContext(c)
	if err != nil {
		return false, err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return false, err
	}

	var owner uint
	err = db.QueryRow("SELECT a.userId FROM accounts a JOIN templates t ON t.accountId = a.id WHERE t.id = $1", template).Scan(&owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"userId":   userID,
			"template": template,
		}).Error("error checking owner of transaction template")
		return false, err
	}

	return owner == userID, nil
}

func templateToDB(transaction Template) *templateDB {
	return &templateDB{
		TemplateName: transaction.TemplateName,
		ID:           transaction.ID,
		Name:         transaction.Name,
		Category:     transaction.Category,
		Amount:       transaction.Amount,
		Note:         transaction.Note,
		AccountID:    transaction.AccountID,
	}
}

func templateFromDB(transaction templateDB) Template {
	return Template{
		TemplateName: transaction.TemplateName,
		Transaction: Transaction{
			ID:        transaction.ID,
			Name:      transaction.Name,
			Category:  transaction.Category,
			Amount:    transaction.Amount,
			Note:      transaction.Note,
			AccountID: transaction.AccountID,
		},
	}
}
