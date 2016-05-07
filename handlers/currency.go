package handlers

import (
	"github.com/emicklei/go-restful"

	"currency"
)

func GetCurrencies(request *restful.Request, response *restful.Response) {
	response.WriteEntity(currency.CodeToName)
}
