package session

import "time"

type Config struct {
	CookieName   string `json:"cookie_name"`   // cookie name of the session
	CookieSecure bool   `json:"cookie_secure"` // is cookie use https?
	CookieExpire int    `json:"cookie_expire"` // cookie expire seconds
	CookieSecret string `json:"cookie_secret"` // cookie secret key

	RememberName   string `json:"remember_name"`   // hashed value for auto login
	RememberExpire int    `json:"remember_expire"` // auto login expire seconds
	RememberSecret string `json:"remember_secret"` // remember cookie secret key

	SessionExpire int `json:"session_expire"` // session expire seconds
}

func (c *Config) CookieExpireSeconds() time.Duration {
	return time.Duration(c.CookieExpire) * time.Second
}

func (c *Config) RememberExpireSeconds() time.Duration {
	return time.Duration(c.RememberExpire) * time.Second
}

func (c *Config) SessionExpireSeconds() time.Duration {
	return time.Duration(c.SessionExpire) * time.Second
}
