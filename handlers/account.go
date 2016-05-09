package handlers

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"account"
)

func GetAccounts(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(string)
	c := appengine.NewContext(request.Request)
	accounts, err := account.Get(c, userId)
	if err != nil {
		log.Errorf(c, "error getting accounts: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(accounts)
}

func NewAccount(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(string)
	c := appengine.NewContext(request.Request)
	acc := new(account.Account)
	err := request.ReadEntity(acc)
	if err != nil {
		log.Errorf(c, "error parsing request to create account: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	acc, err = account.New(c, userId, acc)
	if err == account.InvalidCurrency {
		response.WriteError(http.StatusForbidden, err)
		return
	} else if err != nil {
		log.Errorf(c, "error creating account: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(acc)
}

func UpdateAccount(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(string)
	accountId := request.PathParameter("account-id")
	c := appengine.NewContext(request.Request)
	acc := new(account.Account)
	err := request.ReadEntity(acc)
	if err != nil {
		log.Errorf(c, "error parsing request to update account: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	acc, err = account.Update(c, userId, acc, accountId)
	if err == account.InvalidCurrency {
		response.WriteError(http.StatusForbidden, err)
		return
	} else if err != nil {
		log.Errorf(c, "error creating account: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(acc)
}
