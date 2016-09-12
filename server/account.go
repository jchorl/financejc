package server

import (
	"net/http"

	"github.com/jchorl/financejc/server/account"

	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

func (s server) GetAccounts(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(int)
	accounts, err := account.Get(s.Context(), userId)
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

	acc.User = userId
	acc, err = account.New(s.Context(), acc)
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
	acc := &account.Account{}
	err := request.ReadEntity(acc)
	if err != nil {
		logrus.Errorf("error parsing request to update account: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	acc, err = account.Update(s.Context(), acc)
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
	accountId := request.PathParameter("account-id")
	err := account.Delete(s.Context(), accountId)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeader(http.StatusOK)
}
