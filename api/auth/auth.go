package auth

import (
	"context"
	"net/http"

	"github.com/Sirupsen/logrus"
	"google.golang.org/api/oauth2/v2"

	"github.com/jchorl/financejc/api/user"
)

type Request struct {
	Token string `json:"token"` // provided by google
}

func AuthUser(c context.Context, token string) (user.User, error) {
	googleId, email, err := getGoogleInfoFromToken(token)
	if err != nil {
		return user.User{}, err
	}

	resolved, err := user.FindOrCreateByGoogleId(c, googleId, email)
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
