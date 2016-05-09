package handlers

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"transaction"
)

func GetTransactions(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	accountId := request.PathParameter("account-id")
	transactions, err := transaction.Get(c, accountId)
	if err != nil {
		log.Errorf(c, "error getting transactions: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(transactions)
}

func NewTransaction(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	accountId := request.PathParameter("account-id")
	tr := new(transaction.Transaction)
	err := request.ReadEntity(tr)
	if err != nil {
		log.Errorf(c, "error parsing request to create transaction: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	tr, err = transaction.New(c, accountId, tr)
	if err != nil {
		log.Errorf(c, "error creating transaction: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(tr)
}

func UpdateTransaction(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	transactionId := request.PathParameter("transaction-id")
	tr := new(transaction.Transaction)
	err := request.ReadEntity(tr)
	if err != nil {
		log.Errorf(c, "error parsing request to update transaction: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	tr, err = transaction.Update(c, tr, transactionId)
	if err != nil {
		log.Errorf(c, "error creating transaction: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(tr)
}

func DeleteTransaction(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	transactionId := request.PathParameter("transaction-id")
	err := transaction.Delete(c, transactionId)
	if err != nil {
		log.Errorf(c, "error deleting transaction: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeader(http.StatusOK)
}
