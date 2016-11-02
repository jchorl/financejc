package constants

import (
	"os"
)

// these values can be overridden by environment variables of the same name (DB_DRIVER, etc.)
var (
	JWT_SIGNING_KEY = firstNonEmpty(os.Getenv("JWT_SIGNING_KEY"), "samplejwtsigningkey")
	DB_DRIVER       = firstNonEmpty(os.Getenv("DB_DRIVER"), "postgres")
	DB_ADDRESS      = firstNonEmpty(os.Getenv("DB_ADDRESS"), "postgres://financejc:financejc@financejcdb?sslmode=disable")
	ES_ADDRESS      = firstNonEmpty(os.Getenv("ES_ADDRESS"), "http://financejces:9200")
)

func firstNonEmpty(vals ...string) string {
	for _, val := range vals {
		if val != "" {
			return val
		}
	}

	return ""
}
