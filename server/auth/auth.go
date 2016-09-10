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
	logrus.WithField("Token", token).Debug("attempting to get google ID from token")
	googleId, err := getGoogleIDFromToken(token)
	if err != nil {
		return -1, err
	}

	userId, err := user.FindOrCreateByGoogleId(c, googleId)
	if err != nil {
		return -1, err
	}
	logrus.WithField("User ID", userId).Debug("fetched/created user")

	return userId, nil
}

func getGoogleIDFromToken(token string) (string, error) {
	service, err := oauth2.New(http.DefaultClient)
	tokenInfoCall := service.Tokeninfo()
	tokenInfoCall.IdToken(token)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return "", err
	}

	logrus.Debugf("google responded with tokeninfo: %+v", tokenInfo)
	logrus.WithField("Token Info", tokenInfo).Debug("Google responded")
	return tokenInfo.UserId, nil
}
