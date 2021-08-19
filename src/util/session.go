package util

import (
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/sessions"
	_ "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

type UserSession struct {
	ID     int       `json:"id"`
	UID    string    `json:"uid"`
	Expire time.Time `json:"expire"`
}

var Session *sessions.Session = sessions.NewSession(&sessions.CookieStore{
	Codecs:  nil,
	Options: nil,
}, "memcache")

type SessionStore struct {
	Client *memcache.Client
}

func (c *SessionStore) CreateSSClient(host string) {
	c.Client = memcache.New(host)
	return
}

// GetSession use for getting UserSession by its UserSession.UID
func (c *SessionStore) GetSession(uid string) (*UserSession, error) {
	userS, err := c.Client.Get(uid)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	var tmpUser = &UserSession{}
	err = json.Unmarshal(userS.Value, tmpUser)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return tmpUser, nil
}

// SetSession use for caching UserSession in memory
func (c *SessionStore) SetSession(us *UserSession) bool {
	outV, err := json.Marshal(us)
	err = c.Client.Set(&memcache.Item{
		Key:        us.UID,
		Value:      outV,
		Flags:      0,
		Expiration: 2 * 60 * 60,
	})
	if err != nil {
		log.Error().Msg(err.Error())
		return false
	}

	return true
}
