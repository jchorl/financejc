package financejc

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"server/auth"
	"server/credentials"
)

func init() {
	http.HandleFunc("/auth", authUser)
}

func authUser(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	decoder := json.NewDecoder(r.Body)
	var req auth.Request
	log.Debugf(c, "received request to auth user")
	err := decoder.Decode(&req)
	if err != nil {
		log.Errorf(c, err.Error())
		panic(err)
	}
	userId, err := auth.AuthUser(c, req)
	log.Debugf(c, "authed with userId: %s", userId)
	if err != nil {
		log.Errorf(c, err.Error())
		panic(err)
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["userId"] = userId
	tokenStr, err := token.SignedString([]byte(credentials.JWT_SIGNING_KEY))
	if err != nil {
		log.Errorf(c, err.Error())
		panic(err)
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
