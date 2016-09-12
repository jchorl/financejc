package server

import (
	"net/http"
	"strconv"

	"github.com/jchorl/financejc/server/account"

	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

func (s server) GetAccounts(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(int)
	accounts, err := account.Get(s.ContextWithUser(userId))
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(accounts)
}

func (s server) NewAccount(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(int)
	acc := new(account.Account)
	err := request.ReadEntity(acc)
	if err != nil {
		logrus.Errorf("error parsing request to create account: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	acc, err = account.New(s.ContextWithUser(userId), acc)
	if err == account.InvalidCurrency {
		response.WriteError(http.StatusForbidden, err)
		return
	} else if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(acc)
}

func (s server) UpdateAccount(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(int)
	acc := &account.Account{}
	err := request.ReadEntity(acc)
	if err != nil {
		logrus.Errorf("error parsing request to update account: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	acc, err = account.Update(s.ContextWithUser(userId), acc)
	if err == account.InvalidCurrency {
		response.WriteError(http.StatusForbidden, err)
		return
	} else if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(acc)
}

func (s server) DeleteAccount(request *restful.Request, response *restful.Response) {
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

	err = account.Delete(s.ContextWithUser(userId), accountId)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeader(http.StatusOK)
}
