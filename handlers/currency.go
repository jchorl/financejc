package handlers

import (
	"net/http"

	"currency"
	"util"
)

func GetCurrencies(w http.ResponseWriter, r *http.Request) {
	util.WriteJSONResponse(w, currency.CodeToName)
}
