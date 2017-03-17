package transaction

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/lib/pq"

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

	valid, err := util.UserOwnsAccount(c, accountID)
	if err != nil || !valid {
		return transactions, constants.ErrForbidden
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

// BatchImportTemplates batch imports templates
func BatchImportTemplates(c context.Context, templates []Template) error {
	userID, err := util.UserIDFromContext(c)
	if err != nil || !util.IsUserAdmin(userID) {
		return constants.ErrForbidden
	}

	db, err := util.SQLDBFromContext(c)
	if err != nil {
		return err
	}

	txn, err := db.Begin()
	if err != nil {
		logrus.WithError(err).Error("unable to begin transaction when batch inserting templates")
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn("templates", "id", "templatename", "name", "category", "amount", "note", "accountid"))
	if err != nil {
		logrus.WithError(err).Error("unable to begin copy in when batch inserting templates")
		return err
	}

	for _, template := range templates {
		tdb := templateToDB(template)
		_, err = stmt.Exec(tdb.ID, tdb.TemplateName, tdb.Name, tdb.Category, tdb.Amount, tdb.Note, tdb.AccountID)
		if err != nil {
			logrus.WithError(err).Error("unable to exec template copy when batch inserting templates")
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		logrus.WithError(err).Error("unable to exec batch template copy when batch inserting templates")
		return err
	}

	err = stmt.Close()
	if err != nil {
		logrus.WithError(err).Error("unable to close template copy when batch inserting templates")
		return err
	}

	err = txn.Commit()
	if err != nil {
		logrus.WithError(err).Error("unable to commit template copy when batch inserting templates")
		return err
	}

	return nil
}

// GetAllTemplates queries for all templates
func GetAllTemplates(c context.Context) ([]Template, error) {
	userID, err := util.UserIDFromContext(c)
	if err != nil || !util.IsUserAdmin(userID) {
		return nil, constants.ErrForbidden
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	templates := []Template{}
	rows, err := db.Query("SELECT id, templateName, name, category, amount, note, accountId FROM templates")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to fetch all templates")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var template templateDB
		if err := rows.Scan(&template.ID, &template.TemplateName, &template.Name, &template.Category, &template.Amount, &template.Note, &template.AccountID); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("failed to scan into template")
			return nil, err
		}

		templates = append(templates, templateFromDB(template))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to get all templates from rows")
		return nil, err
	}

	return templates, nil
}

// NewTemplate creates a new template
func NewTemplate(c context.Context, transaction *Template) (*Template, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	valid, err := util.UserOwnsAccount(c, transaction.AccountID)
	if err != nil || !valid {
		return nil, constants.ErrForbidden
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
		return constants.ErrForbidden
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
	userID, err := util.UserIDFromContext(c)
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
