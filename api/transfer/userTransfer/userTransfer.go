package userTransfer

import (
	"bufio"
	"context"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/jchorl/financejc/api/account"
	"github.com/jchorl/financejc/api/transaction"
	"github.com/jchorl/financejc/api/util"
	"github.com/jchorl/financejc/constants"
)

const (
	// currency to use when no currency is present
	defaultCurrency = "USD"

	// states while parsing QIF
	accountState     = "ACCOUNT"
	transactionState = "TRANSACTION"
	noneState        = ""
	optionState      = "OPTION"
)

func round(a float64) int {
	if a < 0 {
		return int(a - 0.5)
	}
	return int(a + 0.5)
}

// Import imports a file for a user
func Import(c context.Context, file io.Reader) error {
	return transferQIF(c, file)
}

func transferQIF(c context.Context, file io.Reader) error {
	userID, err := util.UserIDFromContext(c)
	if err != nil {
		return err
	}

	db, err := util.SQLDBFromContext(c)
	if err != nil {
		return err
	}

	state := noneState
	acc := &account.Account{}
	tr := &transaction.Transaction{}
	uncategorized := make([]*transaction.Transaction, 0)

	scanner := bufio.NewScanner(file)

	tx, err := db.Begin()
	if err != nil {
		logrus.WithError(err).Error("could not begin transaction")
	}

	c = context.WithValue(c, constants.CtxDB, tx)

	for scanner.Scan() {
		line := scanner.Text()

		// skip optional sections
		if strings.HasPrefix(line, "!Option") {
			state = optionState
		} else if strings.HasPrefix(line, "!Clear") {
			state = noneState
		}

		if state == optionState {
			continue
		}

		if line == "!Account" {
			state = accountState
			acc = &account.Account{
				Currency: defaultCurrency,
			}
		} else if strings.HasPrefix(line, "!Type:Cat") {
			state = noneState
		} else if strings.HasPrefix(line, "!Type") {
			state = transactionState
		}

		switch state {
		case accountState:
			switch line[0] {
			case 'N':
				acc.Name = line[1:]
			case '^':
				acc.User = userID
				acc, err = account.New(c, acc)
				if err != nil {
					return err
				}
			}
		case transactionState:
			switch line[0] {
			case 'P':
				tr.Name = line[1:]
			case 'D':
				date, err := time.Parse("2006-01-02", line[1:])
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error":    err,
						"unparsed": line[1:],
					}).Error("could not parse date from QIF")
					return err
				}
				tr.Date = date
			case 'L':
				// could either be a category or an account
				categoryRegex := regexp.MustCompile(`\[(.*)\]`)
				match := categoryRegex.FindStringSubmatch(line[1:])
				if len(match) == 0 {
					tr.Category = strings.Replace(line[1:], "-", "/", -1)
				} else {
					tr.Category = ""
				}
			case 'T':
				amtStr := strings.Replace(line[1:], ",", "", -1)
				amt, err := strconv.ParseFloat(amtStr, 64)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error":    err,
						"unparsed": amtStr,
					}).Error("could not parse amount from QIF")
					return err
				}

				currencyInfo := constants.CurrencyInfo[acc.Currency]
				tr.Amount = round(amt * math.Pow10(currencyInfo.DigitsAfterDecimal))
			case 'M':
				tr.Note = line[1:]
			case '^':
				tr.AccountID = acc.ID
				tr, err = transaction.New(c, tr)
				if err != nil {
					return err
				}

				if tr.Category == "" {
					uncategorized = append(uncategorized, tr)
				}
				tr = &transaction.Transaction{}
			}
		}

		if err := scanner.Err(); err != nil {
			logrus.WithError(err).Error("scanner returned error during import")
			return err
		}
	}

	// now take all the uncategorized transactions and try to pair them up based on date
	for _, tr := range uncategorized {
		if tr.Category == "" {
			// try to find matching transaction
			var tr2 *transaction.Transaction
			found := false
			for _, t := range uncategorized {
				if t.Date == tr.Date && t.Amount == -tr.Amount {
					tr2 = t
					found = true
					break
				}
			}

			if found {
				tr.RelatedTransactionID = tr2.ID
				tr.Category = "Credit Card Payment"
				tr2.RelatedTransactionID = tr.ID
				tr2.Category = "Credit Card Payment"

				_, err = transaction.Update(c, tr)
				if err != nil {
					return err
				}
				_, err = transaction.Update(c, tr2)
				if err != nil {
					return err
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		logrus.WithError(err).Error("could not commit import transaction")
		err2 := tx.Rollback()
		if err2 != nil {
			logrus.WithError(err2).Error("could not rollback import transaction")
		}
		return err
	}

	return nil
}
