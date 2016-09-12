package server

import (
	"net/http"
	"strconv"

	"github.com/jchorl/financejc/server/transaction"

	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

func (s server) GetTransactions(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(int)
	accountIdStr := request.PathParameter("account-id")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":      err,
			"Account ID": accountIdStr,
		}).Error("error parsing account ID to int")
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	next := request.QueryParameter("start")
	transactions, err := transaction.Get(s.ContextWithUser(userId), accountId, next)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	writePaginatedEntity(request, response, transactions)
}

func (s server) NewTransaction(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(int)
	accountIdStr := request.PathParameter("account-id")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":      err,
			"Account ID": accountIdStr,
		}).Error("error parsing account ID to int")
		response.WriteError(http.StatusBadRequest, err)
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
	tr, err = transaction.New(s.ContextWithUser(userId), tr)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(tr)
}

func (s server) UpdateTransaction(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(int)
	tr := &transaction.Transaction{}
	err := request.ReadEntity(tr)
	if err != nil {
		logrus.WithField("Error", err).Error("unable to parse request to update transaction")
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	tr, err = transaction.Update(s.ContextWithUser(userId), tr)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(tr)
}

func (s server) DeleteTransaction(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(int)
	transactionIdStr := request.PathParameter("transaction-id")
	transactionId, err := strconv.Atoi(transactionIdStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":          err,
			"Transaction ID": transactionIdStr,
		}).Error("error parsing transaction ID to int")
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	err = transaction.Delete(s.ContextWithUser(userId), transactionId)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeader(http.StatusOK)
}
