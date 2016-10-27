package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/jchorl/financejc/constants"
)

func GetCurrencies(c echo.Context) error {
	return c.JSON(http.StatusOK, constants.CurrencyInfo)
}
