package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"

	"auth"
	"credentials"
)

func AuthUser(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	var req auth.Request

	err := request.ReadEntity(&req)
	if err != nil {
		log.Errorf(c, "error decoding auth request: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	userId, err := auth.AuthUser(c, req)
	if err != nil {
		log.Errorf(c, "error authenticating user: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	log.Debugf(c, "authed with userId: %s", userId)

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["userId"] = userId
	tokenStr, err := token.SignedString([]byte(credentials.JWT_SIGNING_KEY))
	if err != nil {
		log.Errorf(c, "error getting signed jwt: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	log.Debugf(c, "Token: %s", tokenStr)

	cookie := http.Cookie{
		Name:   "auth",
		Value:  tokenStr,
		MaxAge: 604800,
		// TODO: uncomment for prod
		// Secure:   true,
		// HttpOnly: true,
	}
	http.SetCookie(response, &cookie)
}
