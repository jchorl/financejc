package server

import (
	"github.com/emicklei/go-restful"

	"github.com/jchorl/financejc/constants"
)

func (server) GetCurrencies(request *restful.Request, response *restful.Response) {
	response.WriteEntity(constants.CurrencyCodeToName)
}
