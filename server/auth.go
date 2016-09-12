package server

import (
	"net/http"

	"github.com/jchorl/financejc/constants"
	"github.com/jchorl/financejc/server/auth"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
)

type JWTClaims struct {
	UserId int `json:"userId"`
	jwt.StandardClaims
}

func (s server) CheckAuth(request *restful.Request, response *restful.Response) {
	// if already passed the logged out filter, return 401
	writeError(response, constants.NotLoggedIn)
	return
}

func (s server) AuthUser(request *restful.Request, response *restful.Response) {
	var req auth.Request

	err := request.ReadEntity(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":   err,
			"Request": request,
		}).Error("error decoding auth request")
		writeError(response, constants.BadRequest)
		return
	}

	userId, err := auth.AuthUser(s.Context(), req.Token)
	if err != nil {
		writeError(response, err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{userId, jwt.StandardClaims{}})
	tokenStr, err := token.SignedString([]byte(constants.JWT_SIGNING_KEY))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"Token": token,
		}).Error("error getting signed jwt")
		writeError(response, err)
		return
	}

	cookie := http.Cookie{
		Name:     "auth",
		Value:    tokenStr,
		MaxAge:   604800,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(response, &cookie)
}

func getUserId(unparsedJWT string) (int, error) {
	token, err := jwt.ParseWithClaims(unparsedJWT, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.JWT_SIGNING_KEY), nil
	})
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserId, nil
	} else {
		return -1, err
	}
	return -1, constants.NotLoggedIn
}

// only accept logged in users
func LoggedInFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	cookie, err := request.Request.Cookie("auth")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Error("could not get auth cookie")
		writeError(response, constants.NotLoggedIn)
		return
	}
	userId, err := getUserId(cookie.Value)
	if err != nil && err != constants.NotLoggedIn {
		logrus.WithFields(logrus.Fields{
			"Error":  err,
			"Cookie": cookie,
		}).Error("error parsing jwt")
		writeError(response, constants.NotLoggedIn)
		return
	} else if err == constants.NotLoggedIn {
		writeError(response, constants.NotLoggedIn)
		return
	}
	request.SetAttribute("userId", userId)
	chain.ProcessFilter(request, response)
}

// only accept logged out users
func LoggedOutFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	cookie, err := request.Request.Cookie("auth")
	if err == nil {
		_, err := getUserId(cookie.Value)
		if err == nil {
			return
		}
	}
	chain.ProcessFilter(request, response)
}
