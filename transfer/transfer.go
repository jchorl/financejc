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

	"account"
	"transaction"
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

	scanner := bufio.NewScanner(file)

	return datastore.RunInTransaction(c, func(c context.Context) error {
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

					if amt > 0 {
						tr.Incoming = amt
					} else {
						tr.Outgoing = amt
					}
				case 'M':
					tr.Note = line[1:]
				case '^':
					tr, err = transaction.New(c, acc.Id, tr)
					if err != nil {
						log.Debugf(c, "Error creating account: %s", err.Error())
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
}
