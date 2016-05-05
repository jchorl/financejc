package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"

	"auth"
	"credentials"
	"util"
)

func AuthUser(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	decoder := json.NewDecoder(r.Body)
	var req auth.Request
	log.Debugf(c, "received request to auth user")

	err := decoder.Decode(&req)
	if err != nil {
		log.Errorf(c, "error decoding auth json: %+v", err)
		util.WriteJSONError(w, err)
		return
	}

	userId, err := auth.AuthUser(c, req)
	if err != nil {
		log.Errorf(c, "error authenticating user: %+v", err)
		util.WriteJSONError(w, err)
		return
	}
	log.Debugf(c, "authed with userId: %s", userId)

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["userId"] = userId
	tokenStr, err := token.SignedString([]byte(credentials.JWT_SIGNING_KEY))
	if err != nil {
		log.Errorf(c, "error getting signed jwt: %+v", err)
		util.WriteJSONError(w, err)
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
	http.SetCookie(w, &cookie)
}
