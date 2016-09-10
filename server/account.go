package server

import (
	"net/http"

	"github.com/jchorl/financejc/server/account"

	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

func (s server) GetAccounts(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(string)
	accounts, err := account.Get(s.Context(), userId)
	if err != nil {
		logrus.Errorf("error getting accounts: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(accounts)
}

func (s server) NewAccount(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(string)
	acc := new(account.Account)
	err := request.ReadEntity(acc)
	if err != nil {
		logrus.Errorf("error parsing request to create account: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	acc, err = account.New(s.Context(), userId, acc)
	if err == account.InvalidCurrency {
		response.WriteError(http.StatusForbidden, err)
		return
	} else if err != nil {
		logrus.Errorf("error creating account: %+v", err)
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
		logrus.Errorf("error creating account: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(acc)
}

func (s server) DeleteAccount(request *restful.Request, response *restful.Response) {
	accountId := request.PathParameter("account-id")
	err := account.Delete(s.Context(), accountId)
	if err != nil {
		logrus.Errorf("error deleting account: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeader(http.StatusOK)
}
