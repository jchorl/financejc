package server

import (
	"net/http"
	"strconv"

	"github.com/jchorl/financejc/server/transaction"

	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

func (s server) GetTransactions(request *restful.Request, response *restful.Response) {
	accountId := request.PathParameter("account-id")
	next := request.QueryParameter("start")
	transactions, err := transaction.Get(s.Context(), accountId, next)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	writePaginatedEntity(request, response, transactions)
}

func (s server) NewTransaction(request *restful.Request, response *restful.Response) {
	accountIdStr := request.PathParameter("account-id")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logrus.WithField("Error", err).Error("unable to parse account ID to create transaction")
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	tr := new(transaction.Transaction)
	err = request.ReadEntity(tr)
	if err != nil {
		logrus.WithField("Error", err).Error("unable to parse request to create transaction")
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	tr.Account = accountId
	tr, err = transaction.New(s.Context(), tr)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(tr)
}

func (s server) UpdateTransaction(request *restful.Request, response *restful.Response) {
	tr := &transaction.Transaction{}
	err := request.ReadEntity(tr)
	if err != nil {
		logrus.WithField("Error", err).Error("unable to parse request to update transaction")
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	tr, err = transaction.Update(s.Context(), tr)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(tr)
}

func (s server) DeleteTransaction(request *restful.Request, response *restful.Response) {
	transactionIdStr := request.PathParameter("transaction-id")
	transactionId, err := strconv.Atoi(transactionIdStr)
	if err != nil {
		logrus.WithField("Error", err).Error("unable to parse transaction ID to delete transaction")
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	err = transaction.Delete(s.Context(), transactionId)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeader(http.StatusOK)
}
