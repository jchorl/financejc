package auth

import (
	"context"
	"net/http"

	"github.com/jchorl/financejc/server/user"

	"github.com/Sirupsen/logrus"
	"google.golang.org/api/oauth2/v2"
)

type Request struct {
	Token string `json:"token" description:"Google ID token provided by sign-in flow"`
}

func AuthUser(c context.Context, token string) (int, error) {
	googleId, err := getGoogleIDFromToken(token)
	if err != nil {
		return -1, err
	}

	userId, err := user.FindOrCreateByGoogleId(c, googleId)
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func getGoogleIDFromToken(token string) (string, error) {
	service, err := oauth2.New(http.DefaultClient)
	tokenInfoCall := service.Tokeninfo()
	tokenInfoCall.IdToken(token)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"Token": token,
		}).Error("error calling google to validate token")
		return "", err
	}

	return tokenInfo.UserId, nil
}
