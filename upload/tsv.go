package upload

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TSVParsedTransaction struct {
	Date     string
	Account  string
	Payee    string
	Note     string
	Category string
	Amount   string
}

func (t TSVParsedTransaction) GetName() (string, error) {
	return t.Payee, nil
}

func (t TSVParsedTransaction) GetAccount() (string, error) {
	return t.Account, nil
}

func (t TSVParsedTransaction) GetDate() (time.Time, error) {
	return time.Parse("2006-01-02", t.Date)
}

func (t TSVParsedTransaction) GetCategory() (string, error) {
	return t.Category, nil
}

func (t TSVParsedTransaction) GetIncoming() (float64, error) {
	amt, err := t.parseAmount()
	if err != nil {
		return 0, err
	}

	if amt > 0 {
		return amt, nil
	}
	return 0, nil
}

func (t TSVParsedTransaction) GetOutgoing() (float64, error) {
	amt, err := t.parseAmount()
	if err != nil {
		return 0, err
	}

	if amt < 0 {
		return amt, nil
	}
	return 0, nil
}

func (t TSVParsedTransaction) GetNote() (string, error) {
	return t.Note, nil
}

func (t TSVParsedTransaction) parseAmount() (float64, error) {
	// US amounts have a |
	split := strings.Split(t.Amount, " ")

	// starts with either $ or -$
	str := strings.Replace(split[0], "$", "", -1)

	// remove any commas
	str = strings.Replace(str, ",", "", -1)
	f64, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	return f64, nil
}

func (t TSVParsedTransaction) GetCurrency() (string, error) {
	if strings.Contains(t.Amount, "US") {
		return "USD", nil
	}
	return "CAD", nil
}

func TSVUpload(file *os.File) ([]ParsedTransaction, error) {
	parsedTransactions := make([]ParsedTransaction, 0)
	scanner := bufio.NewScanner(file)
	dateRegex, err := regexp.Compile(`20\d\d-\d\d-\d\d`)
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "\t")

		// can detect real transaction by checking if the first field is a date
		realTransaction := dateRegex.MatchString(split[0])
		if !realTransaction {
			continue
		}
		transaction := TSVParsedTransaction{
			split[0],
			split[1],
			split[3],
			split[4],
			split[5],
			split[6],
		}
		parsedTransactions = append(parsedTransactions, transaction)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return parsedTransactions, nil
}
