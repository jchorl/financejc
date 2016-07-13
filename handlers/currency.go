package handlers

import (
	"github.com/emicklei/go-restful"

	"github.com/jchorl/financejc/currency"
)

func GetCurrencies(request *restful.Request, response *restful.Response) {
	response.WriteEntity(currency.CodeToName)
}
