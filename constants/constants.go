package constants

import (
	"errors"
)

const (
	CTX_DB   = "database"
	CTX_USER = "user"
)

var (
	Forbidden = errors.New("User does not have permission to access this resource.")
)

var CurrencyCodeToName = map[string]string{
	"USD": "United States dollar",
	"EUR": "Euro",
	"JPY": "Japanese yen",
	"GBP": "Pound sterling",
	"AUD": "Australian dollar",
	"CHF": "Swiss franc",
	"CAD": "Canadian dollar",
	"MXN": "Mexican peso",
	"CNY": "Chinese yuan",
	"NZD": "New Zealand dollar",
	"SEK": "Swedish krona",
	"RUB": "Russian ruble",
	"HKD": "Hong Kong dollar",
	"NOK": "Norwegian krone",
	"SGD": "Singapore dollar",
	"TRY": "Turkish lira",
	"KRW": "South Korean won",
	"ZAR": "South African rand",
	"BRL": "Brazilian real",
	"INR": "Indian rupee",
}
