package transaction

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/lib/pq"

	"github.com/jchorl/financejc/api/util"
	"github.com/jchorl/financejc/constants"
	"gopkg.in/olivere/elastic.v5"
)

const (
	esType        = "transaction"
	limitPerQuery = 25
)

// Transactions is a paginated list of transactions, with a link to the next page
type Transactions struct {
	NextLink     string
	Transactions []Transaction `json:"transactions"`
}

// Transaction is a transaction
type Transaction struct {
	ID                   int       `json:"id,omitempty"`
	Name                 string    `json:"name"`
	Date                 time.Time `json:"date"`
	Category             string    `json:"category"`
	Amount               int       `json:"amount"`
	Note                 string    `json:"note"`
	RelatedTransactionID int       `json:"relatedTransactionId,omitempty"`
	AccountID            int       `json:"accountId"`
}

// Query holds params to query transactions by a specific field/value pair
type Query struct {
	Field     string `json:"field"`
	Value     string `json:"value"`
	AccountID int    `json:"accountId"`
}

type transactionDB struct {
	ID                   int
	Name                 string
	Occurred             time.Time
	Category             sql.NullString
	Amount               int
	Note                 sql.NullString
	RelatedTransactionID sql.NullInt64
	AccountID            int
}

type transactionES struct {
	ID                   int       `json:"id,omitempty"`
	Name                 string    `json:"name"`
	Date                 time.Time `json:"date"`
	Category             string    `json:"category"`
	Amount               int       `json:"amount"`
	Note                 string    `json:"note"`
	RelatedTransactionID int       `json:"relatedTransactionId,omitempty"`
	AccountID            int       `json:"accountId"`
	UserID               uint      `json:"userId"`
}

type nextPageParams struct {
	Reference time.Time
	Offset    int
}

// Next returns a link to query for the next page
func (t Transactions) Next() string {
	return t.NextLink
}

// Values returns the actual transactions for the current page
func (t Transactions) Values() (ret []interface{}) {
	for _, tr := range t.Transactions {
		ret = append(ret, tr)
	}

	return ret
}

// InitES initializes elasticsearch with proper analysis config
func InitES(es *elastic.Client) error {
	// largely based on https://qbox.io/blog/multi-field-partial-word-autocomplete-in-elasticsearch-using-ngrams
	resp, err := es.CreateIndex(constants.ESIndex).BodyJson(
		map[string]interface{}{
			"settings": map[string]interface{}{
				"analysis": map[string]interface{}{
					"filter": map[string]interface{}{
						"autocomplete_filter": map[string]interface{}{
							"type":     "edge_ngram",
							"max_gram": 10,
							"token_chars": [...]string{
								"letter",
								"digit",
								"punctuation",
								"symbol",
							},
						},
					},
					"analyzer": map[string]interface{}{
						"autocomplete_analyzer": map[string]interface{}{
							"type":      "custom",
							"tokenizer": "whitespace",
							"filter": [...]string{
								"lowercase",
								"asciifolding",
								"autocomplete_filter",
							},
						},
						"whitespace_analyzer": map[string]interface{}{
							"type":      "custom",
							"tokenizer": "whitespace",
							"filter": [...]string{
								"lowercase",
								"asciifolding",
							},
						},
					},
				},
			},
			"mappings": map[string]interface{}{
				esType: map[string]interface{}{
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":  "integer",
							"index": false,
						},
						"name": map[string]interface{}{
							"type":            "text",
							"analyzer":        "autocomplete_analyzer",
							"search_analyzer": "whitespace_analyzer",
							"fields": map[string]interface{}{
								"raw": map[string]interface{}{
									"type":  "keyword",
									"index": "not_analyzed",
								},
							},
						},
						"date": map[string]string{
							"type": "date",
						},
						"category": map[string]interface{}{
							"type":            "text",
							"analyzer":        "autocomplete_analyzer",
							"search_analyzer": "whitespace_analyzer",
							"fields": map[string]interface{}{
								"raw": map[string]interface{}{
									"type":  "keyword",
									"index": "not_analyzed",
								},
							},
						},
						"amount": map[string]interface{}{
							"type":  "integer",
							"index": false,
						},
						"note": map[string]interface{}{
							"type":            "text",
							"analyzer":        "autocomplete_analyzer",
							"search_analyzer": "whitespace_analyzer",
							"fields": map[string]interface{}{
								"raw": map[string]interface{}{
									"type":  "keyword",
									"index": "not_analyzed",
								},
							},
						},
						"relatedTransactionId": map[string]interface{}{
							"type":  "integer",
							"index": false,
						},
						"accountId": map[string]string{
							"type": "integer",
						},
						"userId": map[string]string{
							"type": "integer",
						},
					},
				},
			},
		},
	).Do(context.Background())
	if err != nil {
		logrus.WithError(err).Error("unable to submit settings to elasticsearch for indexing transactions")
		return err
	} else if !resp.Acknowledged {
		logrus.Error("elasticsearch did not acknowledge creating an index")
		return errors.New("elasticsearch did not acknowledge creating an index")
	}

	return nil
}

// Get fetches transactions for a given account and page parameters
func Get(c context.Context, accountID int, previousNextPageEncoded string) (Transactions, error) {
	valid, err := util.UserOwnsAccount(c, accountID)
	if err != nil || !valid {
		return Transactions{}, constants.ErrForbidden
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return Transactions{}, err
	}

	transactions := Transactions{}

	nextPage := nextPageParams{}
	var rows *sql.Rows
	if previousNextPageEncoded != "" {
		nextPage, err = decodeNextPage(previousNextPageEncoded)
		if err != nil {
			return Transactions{}, err
		}

		rows, err = db.Query("SELECT id, name, occurred, category, amount, note, related_transaction_id, account_id FROM transactions WHERE account_id = $1 AND occurred < $2 ORDER BY occurred DESC, id LIMIT $3 OFFSET $4", accountID, nextPage.Reference, limitPerQuery, nextPage.Offset)
	} else {
		rows, err = db.Query("SELECT id, name, occurred, category, amount, note, related_transaction_id, account_id FROM transactions WHERE account_id = $1 ORDER BY occurred DESC, id LIMIT $2", accountID, limitPerQuery)
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountID,
			"next":      previousNextPageEncoded,
		}).Error("failed to fetch transactions")
		return Transactions{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction transactionDB
		if err := rows.Scan(&transaction.ID, &transaction.Name, &transaction.Occurred, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.RelatedTransactionID, &transaction.AccountID); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"accountId": accountID,
				"next":      previousNextPageEncoded,
			}).Error("failed to scan into transaction")
			return Transactions{}, err
		}

		transactions.Transactions = append(transactions.Transactions, fromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountID,
			"next":      previousNextPageEncoded,
		}).Error("failed to get transactions from rows")
		return Transactions{}, err
	}

	if len(transactions.Transactions) == limitPerQuery {
		// either setting to limitPerQuery (no prev nextPage) or bumping (prev nextPage)
		nextPage.Offset += limitPerQuery

		// if there is no previous nextPage, set the reference
		if previousNextPageEncoded == "" {
			nextPage.Reference = transactions.Transactions[0].Date
		}
		next, err := encodeNextPage(nextPage)
		if err != nil {
			return Transactions{}, err
		}

		transactions.NextLink = next
	}

	return transactions, nil
}

// Summary returns all transactions for a user since a given timestamp
func Summary(ctx context.Context, since time.Time) ([]Transaction, error) {
	userID, err := util.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	db, err := util.DBFromContext(ctx)
	if err != nil {
		return nil, err
	}

	transactions := []Transaction{}
	rows, err := db.Query("SELECT t.id, t.name, t.occurred, t.category, t.amount, t.note, t.related_transaction_id, t.account_id FROM transactions t JOIN accounts a ON t.account_id = a.id WHERE a.user_id = $1 AND t.occurred >= $2 AND t.occurred <= CURRENT_DATE ORDER BY t.occurred DESC", userID, since)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"userID": userID,
			"since":  since,
		}).Error("failed to fetch transactions to find summary")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction transactionDB
		if err := rows.Scan(&transaction.ID, &transaction.Name, &transaction.Occurred, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.RelatedTransactionID, &transaction.AccountID); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err,
				"userID": userID,
				"since":  since,
			}).Error("failed to scan into transaction when finding summary")
			return nil, err
		}

		transactions = append(transactions, fromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"userID": userID,
			"since":  since,
		}).Error("failed to get transactions from rows when finding summary")
		return nil, err
	}

	return transactions, nil
}

// BatchImport batch imports transactions
func BatchImport(c context.Context, transactions []Transaction) error {
	if !util.IsAdminRequest(c) {
		return constants.ErrForbidden
	}

	db, err := util.SQLDBFromContext(c)
	if err != nil {
		return err
	}

	txn, err := db.Begin()
	if err != nil {
		logrus.WithError(err).Error("unable to begin transaction when batch inserting transactions")
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn("transactions", "id", "name", "occurred", "category", "amount", "note", "related_transaction_id", "account_id"))
	if err != nil {
		logrus.WithError(err).Error("unable to begin copy in when batch inserting transactions")
		return err
	}

	for _, transaction := range transactions {
		tdb := toDB(transaction)
		_, err = stmt.Exec(tdb.ID, tdb.Name, tdb.Occurred, tdb.Category, tdb.Amount, tdb.Note, tdb.RelatedTransactionID, tdb.AccountID)
		if err != nil {
			logrus.WithError(err).Error("unable to exec transaction copy when batch inserting transactions")
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		logrus.WithError(err).Error("unable to exec batch transaction copy when batch inserting transactions")
		return err
	}

	err = stmt.Close()
	if err != nil {
		logrus.WithError(err).Error("unable to close transaction copy when batch inserting transactions")
		return err
	}

	err = txn.Commit()
	if err != nil {
		logrus.WithError(err).Error("unable to commit transaction copy when batch inserting transactions")
		return err
	}

	return nil
}

// GetAll queries for all transactions
func GetAll(c context.Context) ([]Transaction, error) {
	if !util.IsAdminRequest(c) {
		return nil, constants.ErrForbidden
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	transactions := []Transaction{}
	rows, err := db.Query("SELECT id, name, occurred, category, amount, note, related_transaction_id, account_id FROM transactions")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to fetch all transactions")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction transactionDB
		if err := rows.Scan(&transaction.ID, &transaction.Name, &transaction.Occurred, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.RelatedTransactionID, &transaction.AccountID); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("failed to scan into transaction")
			return nil, err
		}

		transactions = append(transactions, fromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to get all transactions from rows")
		return nil, err
	}

	return transactions, nil
}

// SearchES does a general search over all fields in ES
func SearchES(ctx context.Context, value string) ([]Transaction, error) {
	userID, err := util.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	es, err := util.ESFromContext(ctx)
	if err != nil {
		return nil, err
	}

	searchResult, err := es.Search().Index(constants.ESIndex).Query(
		elastic.NewBoolQuery().
			Filter(elastic.NewTermQuery("userId", userID)).
			Must(elastic.NewMatchQuery("_all", value).Operator("and").Fuzziness("AUTO")),
	).
		Sort("date", false).
		Size(50).
		Do(context.Background())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"value":   value,
			"context": ctx,
		}).Error("error executing ES search")
		return nil, err
	}

	results := []Transaction{}
	for _, hit := range searchResult.Hits.Hits {
		tr := transactionES{}
		err = json.Unmarshal(*hit.Source, &tr)
		if err != nil {
			logrus.WithError(err).Error("unable to unmarshal json returned from es query")
		}
		results = append(results, fromES(tr))
	}

	return results, nil
}

// QueryES queries elasticsearch given query params
func QueryES(ctx context.Context, query Query) ([]Transaction, error) {
	userID, err := util.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	valid, err := util.UserOwnsAccount(ctx, query.AccountID)
	if err != nil || !valid {
		return nil, constants.ErrForbidden
	}

	es, err := util.ESFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if query.Field != "name" && query.Field != "category" {
		logrus.WithField("query.Field", query.Field).Error("querying for unsupported field")
		return nil, constants.ErrBadRequest
	}

	searchResult, err := es.Search().Index(constants.ESIndex).Query(
		elastic.NewBoolQuery().
			Filter(elastic.NewTermQuery("userId", userID)).
			Must(elastic.NewMatchQuery(query.Field, query.Value).Operator("and").Fuzziness("AUTO")).
			Should(elastic.NewTermQuery("accountId", query.AccountID)),
	).
		Aggregation("top_agg",
			elastic.NewTermsAggregation().
				Field(query.Field+".raw").
				Size(10).SubAggregation("top_agg_hits", elastic.NewTopHitsAggregation().Size(1)),
		).
		Sort("date", false).
		Size(50).
		Do(context.Background())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"query":   query,
			"context": ctx,
		}).Error("error executing ES query by field")
		return nil, err
	}

	results := []Transaction{}
	agg, found := searchResult.Aggregations.Terms("top_agg")
	if !found {
		logrus.WithFields(logrus.Fields{
			"searchResult": searchResult,
		}).Error("top_agg aggregation not found")
		return nil, errors.New("top_agg aggregation not found")
	}

	for _, transactionBucket := range agg.Buckets {
		topHits, found := transactionBucket.TopHits("top_agg_hits")
		if !found {
			logrus.WithFields(logrus.Fields{
				"searchResult":      searchResult,
				"transactionBucket": transactionBucket,
			}).Error("top_agg_hits subaggregation not found")
			return nil, errors.New("top_agg_hits subaggregation not found")
		}
		if topHits.Hits == nil {
			logrus.WithField("topHits", topHits).Error("topHits.Hits should not be nil")
			return nil, errors.New("topHits.Hits should not be nil")
		}

		for _, hit := range topHits.Hits.Hits {
			tr := transactionES{}
			err = json.Unmarshal(*hit.Source, &tr)
			if err != nil {
				logrus.WithError(err).Error("unable to unmarshal json returned from es query")
			}
			results = append(results, fromES(tr))
		}
	}

	return results, nil
}

// New creates a new transaction for a user
func New(ctx context.Context, transaction *Transaction) (*Transaction, error) {
	valid, err := util.UserOwnsAccount(ctx, transaction.AccountID)
	if err != nil || !valid {
		return nil, constants.ErrForbidden
	}

	return newWithoutVerifyingAccountOwnership(ctx, transaction)
}

// newWithoutVerifyingAccountOwnership creates a new transaction without checking that the context has the owner
// of the account. This is useful when generating a transaction on behalf of a user, e.g. recurringTransaction
func newWithoutVerifyingAccountOwnership(ctx context.Context, transaction *Transaction) (*Transaction, error) {
	db, err := util.DBFromContext(ctx)
	if err != nil {
		return nil, err
	}

	tdb := toDB(*transaction)
	var id int
	err = db.QueryRow("INSERT INTO transactions(name, occurred, category, amount, note, related_transaction_id, account_id) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id", tdb.Name, tdb.Occurred, tdb.Category, tdb.Amount, tdb.Note, tdb.RelatedTransactionID, tdb.AccountID).Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":         err,
			"transactionDB": tdb,
			"transaction":   transaction,
		}).Errorf("failed to insert transaction row")
		return nil, err
	}

	transaction.ID = id

	es, err := util.ESFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := util.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, err = es.Index().
		Index(constants.ESIndex).
		Type(esType).
		Id(strconv.Itoa(transaction.ID)).
		BodyJson(toES(transaction, userID)).
		Do(context.Background())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":       err,
			"transaction": transaction,
		}).Error("failed to insert transaction into elasticsearch")
		return nil, err
	}

	return transaction, nil
}

// Update updates a transaction
func Update(ctx context.Context, transaction *Transaction) (*Transaction, error) {
	db, err := util.DBFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check account ownership instead of transaction in case transactions can be moved between accounts in the future
	valid, err := util.UserOwnsAccount(ctx, transaction.AccountID)
	if err != nil || !valid {
		return nil, constants.ErrForbidden
	}

	tdb := toDB(*transaction)
	_, err = db.Exec("UPDATE transactions SET name = $1, occurred = $2, category = $3, amount = $4, note = $5, related_transaction_id = $6 WHERE id = $7", tdb.Name, tdb.Occurred, tdb.Category, tdb.Amount, tdb.Note, tdb.RelatedTransactionID, tdb.ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":         err,
			"transactionDB": tdb,
			"transaction":   transaction,
		}).Errorf("failed to update transaction row")
		return nil, err
	}

	es, err := util.ESFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := util.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// indexing a doc with the same id will replace and bump the version number
	_, err = es.Index().
		Index(constants.ESIndex).
		Type(esType).
		Id(strconv.Itoa(transaction.ID)).
		BodyJson(toES(transaction, userID)).
		Do(context.Background())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":       err,
			"transaction": transaction,
		}).Error("failed to update transaction into elasticsearch")
		return nil, err
	}

	return transaction, nil
}

// Delete deletes a transaction
func Delete(ctx context.Context, transactionID int) error {
	db, err := util.DBFromContext(ctx)
	if err != nil {
		return err
	}

	valid, err := util.UserOwnsTransaction(ctx, transactionID)
	if err != nil || !valid {
		return constants.ErrForbidden
	}

	_, err = db.Exec("DELETE FROM transactions WHERE id = $1", transactionID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":         err,
			"transactionID": transactionID,
		}).Errorf("could not delete transaction")
		return err
	}

	es, err := util.ESFromContext(ctx)
	if err != nil {
		return err
	}

	_, err = es.Delete().
		Index(constants.ESIndex).
		Type(esType).
		Id(strconv.Itoa(transactionID)).
		Do(context.Background())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":         err,
			"transactionId": transactionID,
		}).Error("failed to delete transaction in elasticsearch")
		return err
	}

	return nil
}

// PushAllToES pushes all transactions to elasticsearch
func PushAllToES(c context.Context) error {
	if !util.IsAdminRequest(c) {
		return constants.ErrForbidden
	}

	es, err := util.ESFromContext(c)
	if err != nil {
		return err
	}

	// clear out ES
	_, err = es.DeleteIndex(constants.ESIndex).
		Do(context.Background())
	if err != nil {
		logrus.WithError(err).Error("failed to delete all transactions from elasticsearch")
		return err
	}

	err = InitES(es)
	if err != nil {
		return err
	}

	// get all transactions
	db, err := util.DBFromContext(c)
	if err != nil {
		return err
	}

	esBulkReq := es.Bulk().Index(constants.ESIndex).Type(esType)

	rows, err := db.Query("SELECT t.id, t.name, t.occurred, t.category, t.amount, t.note, t.related_transaction_id, t.account_id, a.user_id FROM transactions t JOIN accounts a on t.account_id = a.id")
	if err != nil {
		logrus.WithError(err).Error("failed to fetch all transactions")
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction transactionDB
		var userID uint
		if err := rows.Scan(&transaction.ID, &transaction.Name, &transaction.Occurred, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.RelatedTransactionID, &transaction.AccountID, &userID); err != nil {
			logrus.WithError(err).Error("failed to scan into transaction")
			return err
		}

		parsed := fromDB(transaction)
		esBulkReq.Add(
			elastic.NewBulkIndexRequest().Id(strconv.Itoa(transaction.ID)).Doc(toES(&parsed, userID)),
		)
	}
	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error("failed to get transactions from rows")
		return err
	}

	_, err = esBulkReq.Do(context.Background())
	if err != nil {
		logrus.WithError(err).Error("failed to bulk post all transactions to es")
		return err
	}

	return nil
}

func encodeNextPage(decoded nextPageParams) (string, error) {
	bts, err := json.Marshal(decoded)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"decoded": decoded,
		}).Error("could not encode next page parameter")
		return "", err
	}

	return string(bts), nil
}

func decodeNextPage(encoded string) (nextPageParams, error) {
	var decoded nextPageParams
	err := json.Unmarshal([]byte(encoded), &decoded)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"encoded": encoded,
		}).Error("could not decode next page parameter")
		return nextPageParams{}, err
	}

	return decoded, nil
}

func toDB(transaction Transaction) *transactionDB {
	return &transactionDB{
		ID:                   transaction.ID,
		Name:                 transaction.Name,
		Occurred:             transaction.Date,
		Category:             util.ToNullStringNonEmpty(transaction.Category),
		Amount:               transaction.Amount,
		Note:                 util.ToNullStringNonEmpty(transaction.Note),
		RelatedTransactionID: util.ToNullIntNonZero(transaction.RelatedTransactionID),
		AccountID:            transaction.AccountID,
	}
}

func fromDB(transaction transactionDB) Transaction {
	return Transaction{
		ID:                   transaction.ID,
		Name:                 transaction.Name,
		Date:                 transaction.Occurred,
		Category:             util.FromNullStringNonEmpty(transaction.Category),
		Amount:               transaction.Amount,
		Note:                 util.FromNullStringNonEmpty(transaction.Note),
		RelatedTransactionID: util.FromNullIntNonZero(transaction.RelatedTransactionID),
		AccountID:            transaction.AccountID,
	}
}

func toES(transaction *Transaction, userID uint) transactionES {
	return transactionES{
		ID:                   transaction.ID,
		Name:                 transaction.Name,
		Date:                 transaction.Date,
		Category:             transaction.Category,
		Amount:               transaction.Amount,
		Note:                 transaction.Note,
		RelatedTransactionID: transaction.RelatedTransactionID,
		AccountID:            transaction.AccountID,
		UserID:               userID,
	}
}

func fromES(transaction transactionES) Transaction {
	return Transaction{
		ID:                   transaction.ID,
		Name:                 transaction.Name,
		Date:                 transaction.Date,
		Category:             transaction.Category,
		Amount:               transaction.Amount,
		Note:                 transaction.Note,
		RelatedTransactionID: transaction.RelatedTransactionID,
		AccountID:            transaction.AccountID,
	}
}
