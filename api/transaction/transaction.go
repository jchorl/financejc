package transaction

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"

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
						"note": map[string]string{
							"type": "text",
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
func Get(c context.Context, accountID int, nextEncoded string) (Transactions, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return Transactions{}, err
	}

	valid, err := userOwnsAccount(c, accountID)
	if err != nil || !valid {
		return Transactions{}, constants.ErrForbidden
	}

	transactions := Transactions{}

	reference := time.Now().Add(time.Hour * time.Duration(24))
	offset := 0
	if nextEncoded != "" {
		decoded, err := decodeNextPage(nextEncoded)
		if err != nil {
			return Transactions{}, err
		}

		reference, offset = decoded.Reference, decoded.Offset
	}

	logrus.WithFields(logrus.Fields{
		"accountId":     accountID,
		"reference":     reference,
		"limitPerQuery": limitPerQuery,
		"offset":        offset,
	}).Info("about to query")
	rows, err := db.Query("SELECT id, name, occurred, category, amount, note, relatedTransactionId, accountId FROM transactions WHERE accountId = $1 AND occurred < $2 ORDER BY occurred DESC, id LIMIT $3 OFFSET $4", accountID, reference, limitPerQuery, offset)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountID,
			"next":      nextEncoded,
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
				"next":      nextEncoded,
			}).Error("failed to scan into transaction")
			return Transactions{}, err
		}

		transactions.Transactions = append(transactions.Transactions, fromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountID,
			"next":      nextEncoded,
		}).Error("failed to get transactions from rows")
		return Transactions{}, err
	}

	if len(transactions.Transactions) == limitPerQuery {
		next, err := encodeNextPage(nextPageParams{reference, offset + limitPerQuery})
		if err != nil {
			return Transactions{}, err
		}

		transactions.NextLink = next
	}

	return transactions, nil
}

// GetFuture gets all transactions for an account after a given reference time
func GetFuture(c context.Context, accountID int, reference *time.Time) ([]Transaction, error) {
	db, err := util.DBFromContext(c)
	if err != nil {
		return nil, err
	}

	valid, err := userOwnsAccount(c, accountID)
	if err != nil || !valid {
		return nil, constants.ErrForbidden
	}

	transactions := []Transaction{}

	if reference == nil {
		now := time.Now()
		reference = &now
	}
	rows, err := db.Query("SELECT id, name, occurred, category, amount, note, relatedTransactionId, accountId FROM transactions WHERE accountId = $1 AND occurred > $2 ORDER BY occurred DESC, id", accountID, reference)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountID,
		}).Error("failed to fetch future transactions")
		return transactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction transactionDB
		if err := rows.Scan(&transaction.ID, &transaction.Name, &transaction.Occurred, &transaction.Category, &transaction.Amount, &transaction.Note, &transaction.RelatedTransactionID, &transaction.AccountID); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"accountId": accountID,
			}).Error("failed to scan into transaction for future fetch")
			return transactions, err
		}

		transactions = append(transactions, fromDB(transaction))
	}
	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"accountId": accountID,
		}).Error("failed to get transactions from rows for future fetch")
		return transactions, err
	}

	return transactions, nil
}

// QueryES queries elasticsearch given query params
func QueryES(ctx context.Context, query Query) ([]Transaction, error) {
	userID, err := util.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	valid, err := userOwnsAccount(ctx, query.AccountID)
	if err != nil || !valid {
		return nil, constants.ErrForbidden
	}

	if query.Field != "name" && query.Field != "category" {
		logrus.WithField("query.Field", query.Field).Error("querying for unsupported field")
		return nil, constants.ErrBadRequest
	}

	es, err := util.ESFromContext(ctx)
	if err != nil {
		return nil, err
	}

	searchResult, err := es.Search().
		Index(constants.ESIndex).
		Query(
			elastic.NewBoolQuery().
				Filter(elastic.NewTermQuery("userId", userID)).
				Must(elastic.NewMatchQuery(query.Field, query.Value).Operator("and").Fuzziness("AUTO")).
				Should(elastic.NewTermQuery("accountId", query.AccountID))).
		Aggregation("top_agg",
			elastic.NewTermsAggregation().Field(query.Field+".raw").Size(10).SubAggregation(
				"top_agg_hits", elastic.NewTopHitsAggregation().Size(1))).
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
	valid, err := userOwnsAccount(ctx, transaction.AccountID)
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
	err = db.QueryRow("INSERT INTO transactions(name, occurred, category, amount, note, relatedTransactionId, accountId) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id", tdb.Name, tdb.Occurred, tdb.Category, tdb.Amount, tdb.Note, tdb.RelatedTransactionID, tdb.AccountID).Scan(&id)
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

	valid, err := userOwnsAccount(ctx, transaction.AccountID)
	if err != nil || !valid {
		return nil, constants.ErrForbidden
	}

	tdb := toDB(*transaction)
	_, err = db.Exec("UPDATE transactions SET name = $1, occurred = $2, category = $3, amount = $4, note = $5, relatedTransactionId = $6 WHERE id = $7", tdb.Name, tdb.Occurred, tdb.Category, tdb.Amount, tdb.Note, tdb.RelatedTransactionID, tdb.ID)
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

	valid, err := userOwnsTransaction(ctx, transactionID)
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
	userID, err := util.UserIDFromContext(c)
	if err != nil || userID != constants.AdminUID {
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

	rows, err := db.Query("SELECT t.id, t.name, t.occurred, t.category, t.amount, t.note, t.relatedTransactionId, t.accountId, a.userId FROM transactions t JOIN accounts a on t.accountId = a.id")
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

func userOwnsAccount(c context.Context, account int) (bool, error) {
	userID, err := util.UserIDFromContext(c)
	if err != nil {
		return false, err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return false, err
	}

	var owner uint
	err = db.QueryRow("SELECT userId FROM accounts WHERE id = $1", account).Scan(&owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"userId":  userID,
			"account": account,
		}).Error("error checking owner of account")
		return false, err
	}

	return owner == userID, nil
}

func userOwnsTransaction(c context.Context, transaction int) (bool, error) {
	userID, err := util.UserIDFromContext(c)
	if err != nil {
		return false, err
	}

	db, err := util.DBFromContext(c)
	if err != nil {
		return false, err
	}

	var owner uint
	err = db.QueryRow("SELECT a.userId FROM accounts a JOIN transactions t ON t.accountId = a.id WHERE t.id = $1", transaction).Scan(&owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":       err,
			"userId":      userID,
			"transaction": transaction,
		}).Error("error checking owner of transaction")
		return false, err
	}

	return owner == userID, nil
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
