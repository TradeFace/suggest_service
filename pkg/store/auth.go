package store

import (
	"errors"
	"time"

	"github.com/tradeface/jwt_service/pkg/authorization"
)

const CACHE_LASTSEEN_MIN = 15

type AuthStore struct {
	authUser map[string]*authorization.AuthUser
}

func NewAuthStore() *AuthStore {

	authStore := &AuthStore{
		authUser: make(map[string]*authorization.AuthUser, 0),
	}

	done := make(chan bool)
	authStore.forever(done)

	return authStore
}

func (ac *AuthStore) GetAuthUser(tokenString string) (authUser *authorization.AuthUser, err error) {

	if authUser, ok := ac.authUser[tokenString]; ok {
		//renew lastseen
		authUser.LastseenExpire = time.Now().Add(CACHE_LASTSEEN_MIN * time.Minute)
		return authUser, nil
	}
	return authUser, errors.New("user not found")
}

func (ac *AuthStore) AddAuthUser(tokenString string, authUser *authorization.AuthUser) (err error) {

	expiresAt, err := authUser.GetClaim("ExpiresAt")
	if err != nil {
		return err
	}
	authUser.TokenExpire = time.Unix(expiresAt.(int64), 0)
	authUser.LastseenExpire = time.Now().Add(CACHE_LASTSEEN_MIN * time.Minute)
	ac.authUser[tokenString] = authUser
	return err
}

func (ac *AuthStore) forever(done <-chan bool) {

	//remove all expired AuthUsers
	//remove all lastseen AuthUsers (we don't want to keep long lived user in mem)
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				for key, val := range ac.authUser {
					if val.TokenExpire.Before(time.Now()) {
						delete(ac.authUser, key)
						continue
					}
					if val.LastseenExpire.Before(time.Now()) {
						delete(ac.authUser, key)
						continue
					}
				}
			}
		}
	}()
}
