package transfer

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"github.com/jchorl/financejc/account"
	"github.com/jchorl/financejc/transaction"
)

const (
	ACCOUNT     = "ACCOUNT"
	TRANSACTION = "TRANSACTION"
	NONE        = ""
	OPTION      = "OPTION"
)

func TransferQIF(c context.Context, userId string, file *os.File) error {
	var err error
	state := NONE
	acc := &account.Account{}
	tr := &transaction.Transaction{}
	uncategorized := make([]*transaction.Transaction, 0)

	scanner := bufio.NewScanner(file)

	err = datastore.RunInTransaction(c, func(c context.Context) error {
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
					Currency: "USD",
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
					acc, err = account.New(c, userId, acc)
					if err != nil {
						log.Debugf(c, "Error creating account: %s", err.Error())
					}
				}
			case TRANSACTION:
				switch line[0] {
				case 'P':
					tr.Name = line[1:]
				case 'D':
					date, err := time.Parse("2006-01-02", line[1:])
					if err != nil {
						log.Debugf(c, "Error parsing date: %s", err.Error())
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
						log.Debugf(c, "Error parsing amount: %s", err.Error())
						return err
					}

					tr.Amount = amt
				case 'M':
					tr.Note = line[1:]
				case '^':
					tr, err = transaction.New(c, acc.Id, tr)
					if err != nil {
						log.Debugf(c, "Error creating account: %s", err.Error())
						return err
					}

					if tr.Category == "" {
						uncategorized = append(uncategorized, tr)
					}
					tr = &transaction.Transaction{}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
		return nil
	}, nil)

	if err != nil {
		log.Debugf(c, "Error executing transaction: %s", err.Error())
		return err
	}

	return datastore.RunInTransaction(c, func(c context.Context) error {
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
					tr.RelatedTransaction = tr2.Id
					tr.Category = "Credit Card Payment"
					tr2.RelatedTransaction = tr.Id
					tr2.Category = "Credit Card Payment"

					_, err = transaction.Update(c, tr, tr.Id)
					if err != nil {
						log.Debugf(c, "Error updating transaction: %s", err.Error())
						return err
					}
					_, err = transaction.Update(c, tr2, tr2.Id)
					if err != nil {
						log.Debugf(c, "Error updating transaction: %s", err.Error())
						return err
					}
				}
			}
		}

		return nil
	}, nil)
}
