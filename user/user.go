package user

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

const dbKey string = "User"

type User struct {
	GoogleId string `datastore:"google",json:"-"`
}

func Key(userId string) (*datastore.Key, error) {
	return datastore.DecodeKey(userId)
}

func FindOrCreateByGoogleId(c context.Context, googleId string) (string, error) {
	q := datastore.NewQuery(dbKey).
		Filter("GoogleId =", googleId).
		KeysOnly().
		Limit(1)
	keys, err := q.GetAll(c, nil)
	if err != nil {
		return "", err
	}
	if len(keys) == 1 {
		return keys[0].Encode(), nil
	}

	user := &User{
		GoogleId: googleId,
	}

	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, dbKey, nil), user)
	if err != nil {
		return "", err
	}
	return key.Encode(), nil
}
