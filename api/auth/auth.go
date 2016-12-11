package auth

import (
	"context"
	"net/http"

	"github.com/Sirupsen/logrus"
	"google.golang.org/api/oauth2/v2"

	"github.com/jchorl/financejc/api/user"
)

// Request represents a login request from a client
type Request struct {
	Token string `json:"token"` // provided by google
}

// LoginByGoogleToken takes a token from Google and returns a User
func LoginByGoogleToken(c context.Context, token string) (user.User, error) {
	googleID, email, err := getGoogleInfoFromToken(token)
	if err != nil {
		return user.User{}, err
	}

	resolved, err := user.FindOrCreateByGoogleId(c, googleID, email)
	if err != nil {
		return user.User{}, err
	}

	return resolved, nil
}

func getGoogleInfoFromToken(token string) (string, string, error) {
	service, err := oauth2.New(http.DefaultClient)
	tokenInfoCall := service.Tokeninfo()
	tokenInfoCall.IdToken(token)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"token": token,
		}).Error("error calling google to validate token")
		return "", "", err
	}

	return tokenInfo.UserId, tokenInfo.Email, nil
}
