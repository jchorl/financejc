package constants

import (
	"os"

	"github.com/jchorl/financejc/server/util"
)

// these values can be overridden by environment variables of the same name (DB_DRIVER, etc.)
var (
	JWT_SIGNING_KEY = util.FirstNonEmpty(os.Getenv("JWT_SIGNING_KEY"), "samplejwtsigningkey")
	DB_DRIVER       = util.FirstNonEmpty(os.Getenv("DB_DRIVER"), "postgres")
	DB_ADDRESS      = util.FirstNonEmpty(os.Getenv("DB_ADDRESS"), "postgres://financejc:financejc@financejcdb?sslmode=disable")
)
