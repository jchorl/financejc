package transfer

import (
	"bufio"
	"context"
	"io/ioutil"
	"math"
	"os"
	"path"
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
	ACCOUNT          = "ACCOUNT"
	TRANSACTION      = "TRANSACTION"
	NONE             = ""
	OPTION           = "OPTION"
	DEFAULT_CURRENCY = "USD"
)

func round(a float64) int {
	if a < 0 {
		return int(a - 0.5)
	}
	return int(a + 0.5)
}

func AutoImport(c context.Context) error {
	files, err := ioutil.ReadDir(constants.IMPORT_PATH)
	if err != nil {
		currDir, err2 := os.Getwd()
		if err2 != nil {
			logrus.WithField("Error", err2).Error("error listing working dir while reporting listing error")
		}
		logrus.WithFields(logrus.Fields{
			"error":            err,
			"importPath":       constants.IMPORT_PATH,
			"currentDirectory": currDir,
		}).Error("error listing files")
		return err
	}
	for _, f := range files {
		// skip gitkeep
		if f.Name() == ".gitkeep" {
			continue
		}

		file, err := os.Open(path.Join(constants.IMPORT_PATH, f.Name()))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":    err,
				"file":     f.Name(),
				"fullPath": path.Join(constants.IMPORT_PATH, f.Name()),
			}).Error("error opening files")
			return err
		}
		defer file.Close()

		err = TransferQIF(c, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func TransferQIF(c context.Context, file *os.File) error {
	userId, err := util.UserIdFromContext(c)
	if err != nil {
		return err
	}

	db, err := util.SQLDBFromContext(c)
	if err != nil {
		return err
	}

	state := NONE
	acc := &account.Account{}
	tr := &transaction.Transaction{}
	uncategorized := make([]*transaction.Transaction, 0)

	scanner := bufio.NewScanner(file)

	tx, err := db.Begin()
	if err != nil {
		logrus.WithError(err).Error("could not begin transaction")
	}

	c = context.WithValue(c, constants.CTX_DB, tx)

	for scanner.Scan() {
		line := scanner.Text()

		// skip optional sections
		if strings.HasPrefix(line, "!Option") {
			state = OPTION
		} else if strings.HasPrefix(line, "!Clear") {
			state = NONE
		}

		if state == OPTION {
			continue
		}

		if line == "!Account" {
			state = ACCOUNT
			acc = &account.Account{
				Currency: DEFAULT_CURRENCY,
			}
		} else if strings.HasPrefix(line, "!Type:Cat") {
			state = NONE
		} else if strings.HasPrefix(line, "!Type") {
			state = TRANSACTION
		}

		switch state {
		case ACCOUNT:
			switch line[0] {
			case 'N':
				acc.Name = line[1:]
			case '^':
				acc.User = userId
				acc, err = account.New(c, acc)
				if err != nil {
					return err
				}
			}
		case TRANSACTION:
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
				tr.AccountId = acc.Id
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
				tr.RelatedTransactionId = tr2.Id
				tr.Category = "Credit Card Payment"
				tr2.RelatedTransactionId = tr.Id
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
