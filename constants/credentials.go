package constants

import (
	"os"
)

// these values can be overridden by environment variables of the same name (DB_DRIVER, etc.)
var (
	// JwtSigningKey is the key used to sign JWTs
	JwtSigningKey = firstNonEmpty(os.Getenv("JWT_SIGNING_KEY"), "samplejwtsigningkey")

	// DbDriver is the driver for the db
	DbDriver = firstNonEmpty(os.Getenv("DB_DRIVER"), "postgres")

	// DbAddress is the address of the db
	DbAddress = firstNonEmpty(os.Getenv("DB_ADDRESS"), "postgres://postgres@financejcdb?sslmode=disable")

	// EsAddress is the address of elasticsearch
	EsAddress = firstNonEmpty(os.Getenv("ES_ADDRESS"), "http://financejces:9200")

	// GcsAccountJSON is a json service account credentials generated by google's api credentials
	GcsAccountJSON = os.Getenv("GCS_ACCOUNT_JSON")
)

func firstNonEmpty(vals ...string) string {
	for _, val := range vals {
		if val != "" {
			return val
		}
	}

	return ""
}
