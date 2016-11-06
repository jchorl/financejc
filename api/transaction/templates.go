package transaction

import (
	"context"

	"github.com/Sirupsen/logrus"

	"github.com/jchorl/financejc/api/util"
	"github.com/jchorl/financejc/constants"
)

type TransactionTemplate Transaction

type transactionTemplateDB struct {
	Id        int
	Name      string
	Category  string
	Amount    int
	Note      string
	AccountId int
}

func GetTemplates(c context.Context, accountId int) ([]TransactionTemplate, error) {
	transactions := []TransactionTemplate{}
	db, err := util.DBFromContext(c)
	if err != nil {
		return transactions, err
	}

	valid, err := userOwnsAccount(c, accountId)
	if err != nil || !valid {
		return transactions, constants.Forbidden
	}

	rows, err := db.Query("SELECT id, name, category, amount, note, accountId FROM transactionTemplates WHERE accountId = $1", accountId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountId,
		}).Error("failed to fetch transaction templates")
		return transactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction transactionTemplateDB
		if err := rows.Scan(&transaction.Id, &transaction.Name, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.AccountId); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"accountId": accountId,
			}).Error("failed to scan into transaction template")
			return transactions, err
		}

		transactions = append(transactions, templateFromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountId,
		}).Error("failed to get transaction templates from rows")
		return transactions, err
	}

	return transactions, nil
}

func NewTemplate(c context.Context, transaction *TransactionTemplate) (*TransactionTemplate, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	valid, err := userOwnsAccount(c, transaction.AccountId)
	if err != nil || !valid {
		return nil, constants.Forbidden
	}

	tdb := templateToDB(*transaction)
	var id int
	err = db.QueryRow("INSERT INTO transactionTemplates(name, category, amount, note, accountId) VALUES($1, $2, $3, $4, $5) RETURNING id", tdb.Name, tdb.Category, tdb.Amount, tdb.Note, tdb.AccountId).Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":                 err,
			"transactionTemplateDB": tdb,
			"transactionTemplate":   transaction,
		}).Errorf("failed to insert transaction template row")
		return nil, err
	}

	transaction.Id = id
	return transaction, nil
}

func UpdateTemplate(c context.Context, transaction *TransactionTemplate) (*TransactionTemplate, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	tdb := templateToDB(*transaction)
	_, err = db.Exec("UPDATE transactionTemplates SET name = $1, category = $2, amount = $3, note = $4, accountId = $5", tdb.Name, tdb.Category, tdb.Amount, tdb.Note, tdb.AccountId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":                 err,
			"transactionTemplateDB": tdb,
			"transactionTemplate":   transaction,
		}).Errorf("failed to update transaction template row")
		return nil, err
	}

	return transaction, nil
}

func DeleteTemplate(ctx context.Context, transactionId int) error {
	db, err := util.DBFromContext(ctx)
	if err != nil {
		return err
	}

	valid, err := userOwnsTransactionTemplate(ctx, transactionId)
	if err != nil || !valid {
		return constants.Forbidden
	}

	_, err = db.Exec("DELETE FROM transactionTemplates WHERE id = $1", transactionId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":                 err,
			"transactionTemplateId": transactionId,
		}).Errorf("could not delete transaction template")
		return err
	}

	return nil
}

func userOwnsTransactionTemplate(c context.Context, transactionTemplate int) (bool, error) {
	userId, err := util.UserIdFromContext(c)
	if err != nil {
		return false, err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return false, err
	}

	var owner uint
	err = db.QueryRow("SELECT a.userId FROM accounts a JOIN transactionTemplates t ON t.accountId = a.id WHERE t.id = $1", transactionTemplate).Scan(&owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":               err,
			"userId":              userId,
			"transactionTemplate": transactionTemplate,
		}).Error("error checking owner of transaction template")
		return false, err
	}

	return owner == userId, nil
}

func templateToDB(transaction TransactionTemplate) *transactionTemplateDB {
	return &transactionTemplateDB{
		Id:        transaction.Id,
		Name:      transaction.Name,
		Category:  transaction.Category,
		Amount:    transaction.Amount,
		Note:      transaction.Note,
		AccountId: transaction.AccountId,
	}
}

func templateFromDB(transaction transactionTemplateDB) TransactionTemplate {
	return TransactionTemplate{
		Id:        transaction.Id,
		Name:      transaction.Name,
		Category:  transaction.Category,
		Amount:    transaction.Amount,
		Note:      transaction.Note,
		AccountId: transaction.AccountId,
	}
}
