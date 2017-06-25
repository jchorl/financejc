package batchTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/jchorl/financejc/api/account"
	"github.com/jchorl/financejc/api/transaction"
	"github.com/jchorl/financejc/api/user"
	"github.com/jchorl/financejc/api/util"
	"github.com/jchorl/financejc/constants"
)

type fjcData struct {
	Users                 []user.User                        `json:"users"`
	Accounts              []account.Account                  `json:"accounts"`
	Transactions          []transaction.Transaction          `json:"transactions"`
	RecurringTransactions []transaction.RecurringTransaction `json:"recurringTransactions"`
	Templates             []transaction.Template             `json:"templates"`
}

// BackupToGCS exports all app data and pushes it to a GCS bucket
func BackupToGCS(c context.Context) error {
	logrus.Debug("starting regular backup to GCS")
	if !util.IsAdminRequest(c) {
		return constants.ErrForbidden
	}

	conf, err := google.JWTConfigFromJSON([]byte(constants.GcsAccountJSON), storage.ScopeReadWrite)
	if err != nil {
		logrus.WithError(err).Error("failed to create jwt config from json")
		return err
	}

	gctx := context.Background()
	client, err := storage.NewClient(
		gctx,
		option.WithTokenSource(conf.TokenSource(gctx)),
	)
	if err != nil {
		logrus.WithError(err).Error("unable to create storage client")
		return err
	}

	backupsBucketExists := false
	bucketIter := client.Buckets(gctx, constants.GoogleProjectID)
	for {
		bucketAttrs, err := bucketIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logrus.WithError(err).Error("error iterating gcs buckets")
			return err
		}

		if bucketAttrs.Name == constants.GcsBackupBucket {
			backupsBucketExists = true
		}
	}

	bucket := client.Bucket(constants.GcsBackupBucket)
	if !backupsBucketExists {
		if err := bucket.Create(gctx, constants.GoogleProjectID, &storage.BucketAttrs{StorageClass: "REGIONAL", Location: "us-central1"}); err != nil {
			logrus.WithError(err).Error("error creating backup bucket")
			return err
		}
	}

	content, err := Export(c)
	if err != nil {
		logrus.WithError(err).Error("failed to generate export for periodic backup")
		return err
	}

	filename := time.Now().Format("20060102T150405")
	obj := bucket.Object(filename)
	w := obj.NewWriter(gctx)
	if _, err := fmt.Fprintf(w, content); err != nil {
		logrus.WithError(err).Error("unable to write to object when creating backup")
		return err
	}
	if err := w.Close(); err != nil {
		logrus.WithError(err).Error("unable to close writer when creating backup")
		return err
	}

	logrus.Debugf("backup finished successfully with filename: %s", filename)
	return nil
}

// Export queries for all data, packages it up and exports it
func Export(c context.Context) (string, error) {
	if !util.IsAdminRequest(c) {
		return "", constants.ErrForbidden
	}

	allData := fjcData{}
	users, err := user.GetAll(c)
	if err != nil {
		return "", err
	}
	allData.Users = users

	accounts, err := account.GetAll(c)
	if err != nil {
		return "", err
	}
	allData.Accounts = accounts

	transactions, err := transaction.GetAll(c)
	if err != nil {
		return "", err
	}
	allData.Transactions = transactions

	templates, err := transaction.GetAllTemplates(c)
	if err != nil {
		return "", err
	}
	allData.Templates = templates

	recurringTransactions, err := transaction.GetAllRecurring(c)
	if err != nil {
		return "", err
	}
	allData.RecurringTransactions = recurringTransactions

	return encode(allData)
}

// Import batch imports the result of an export
func Import(c context.Context, encoded string) error {
	if !util.IsAdminRequest(c) {
		return constants.ErrForbidden
	}

	allData, err := decode(encoded)
	if err != nil {
		return err
	}

	if err = user.BatchImport(c, allData.Users); err != nil {
		return err
	}

	if err = account.BatchImport(c, allData.Accounts); err != nil {
		return err
	}

	if err = transaction.BatchImport(c, allData.Transactions); err != nil {
		return err
	}

	if err = transaction.BatchImportTemplates(c, allData.Templates); err != nil {
		return err
	}

	if err = transaction.BatchImportRecurringTransactions(c, allData.RecurringTransactions); err != nil {
		return err
	}

	// update all the auto-increment sequences
	db, err := util.DBFromContext(c)
	if err != nil {
		return err
	}

	_, err = db.Query(`SELECT setval('users_id_seq', (SELECT MAX(id) from "users"));`)
	if err != nil {
		logrus.WithError(err).Error("unable to update the users sequence")
		return err
	}

	_, err = db.Query(`SELECT setval('accounts_id_seq', (SELECT MAX(id) from "accounts"));`)
	if err != nil {
		logrus.WithError(err).Error("unable to update the accounts sequence")
		return err
	}

	_, err = db.Query(`SELECT setval('transactions_id_seq', (SELECT MAX(id) from "transactions"));`)
	if err != nil {
		logrus.WithError(err).Error("unable to update the transactions sequence")
		return err
	}

	_, err = db.Query(`SELECT setval('recurring_transactions_id_seq', (SELECT MAX(id) from "recurring_transactions"));`)
	if err != nil {
		logrus.WithError(err).Error("unable to update the recurring_transactions sequence")
		return err
	}

	_, err = db.Query(`SELECT setval('templates_id_seq', (SELECT MAX(id) from "templates"));`)
	if err != nil {
		logrus.WithError(err).Error("unable to update the templates sequence")
		return err
	}

	return nil
}

func encode(data fjcData) (string, error) {
	encB, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Error("error encoding all fjc data")
		return "", err
	}

	return string(encB), nil
}

func decode(str string) (fjcData, error) {
	data := fjcData{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		logrus.WithError(err).Error("error decoding all fjc data")
		return data, err
	}

	return data, nil
}
