package session

import (
	"net/http"
	"net/url"
	"time"
)

const (
	MaxRetried = 3
)

// Session delegates to provider of session store
type Session struct {
	config   *Config
	provider Provider
}

// New returns a new session with provider
func New(provider Provider, config *Config) *Session {
	// SecretKey MUST be present or panic!
	if config.CookieSecret == "" {
		panic(ErrEmptySecretKey)
	}

	if config.RememberSecret == "" {
		config.RememberSecret = config.CookieSecret
	}

	return &Session{
		config:   config,
		provider: provider,
	}
}

func (sess *Session) New(w http.ResponseWriter, r *http.Request) (sto Storer, err error) {
	var (
		ok bool
	)

	for i := 0; i < MaxRetried; i++ {
		sid := Helpers.NewSessionID()

		sto, err = sess.provider.New(sid)
		if err == nil {
			ok = true
			break
		}
	}

	// it's ok to ignore errors
	if ok {
		sess.setCookie(w, r, sess.config.CookieName, sto.SessionID(), sess.config.CookieSecret, sess.config.CookieExpire)
	}

	return
}

// Start restores a session from request cookie or creates a new one if absence
func (sess *Session) Start(w http.ResponseWriter, r *http.Request) (sto Storer, err error) {
	cookie, err := sess.getCookie(r, sess.config.CookieName, sess.config.CookieSecret)
	if err != nil {
		if err == ErrNotFound {
			return sess.New(w, r)
		}

		return
	}

	if cookie.IsExpiredBy(sess.config.CookieExpireSeconds()) {
		err = ErrCookieExpired
		return
	}

	// restore session from provider
	sto, err = sess.provider.Restore(cookie.Value)
	if err == ErrNotFound {
		return sess.New(w, r)
	}

	return
}

// Refresh updates request session with new ID
func (sess *Session) Refresh(w http.ResponseWriter, r *http.Request) (sto Storer, err error) {
	cookie, err := sess.getCookie(r, sess.config.CookieName, sess.config.CookieSecret)
	if err != nil {
		return
	}

	if cookie.IsExpiredBy(sess.config.CookieExpireSeconds()) {
		err = ErrCookieExpired
		return
	}

	// restore session from provider
	oldsto, err := sess.provider.Restore(cookie.Value)
	if err != nil {
		return
	}
	defer sess.provider.Destroy(oldsto.SessionID())

	sto, err = sess.New(w, r)
	if err != nil {
		return
	}

	// copy all oldsto data to sto
	value := oldsto.GetValue()
	value.changed = true

	err = sto.SetValue(oldsto.GetValue())

	return
}

// Destroy deletes session of current request
func (sess *Session) Destroy(w http.ResponseWriter, r *http.Request) error {
	cookie, err := sess.getCookie(r, sess.config.CookieName, sess.config.CookieSecret)
	if err != nil {
		if err == ErrNotFound {
			err = nil
		}

		return err
	}

	// ignore provider error is ok!
	sess.provider.Destroy(cookie.Value)

	// force expire client cookie
	sess.delCookie(w, sess.config.CookieName)

	return nil
}

func (sess *Session) setCookie(w http.ResponseWriter, r *http.Request, name, value, secret string, expires int) (cookie *Cookie, err error) {
	cookie = NewCookie(value)
	encoding := NewCookieEncoding(name, secret)

	// secure cookie value of sid
	data, err := encoding.Encrypt(cookie)
	if err != nil {
		return
	}

	tmpcookie := &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(data),
		Path:     "/",
		Secure:   sess.config.CookieSecure,
		HttpOnly: true,
	}

	if expires >= 0 { // unit in seconds
		tmpcookie.MaxAge = expires
	}

	http.SetCookie(w, tmpcookie)
	r.AddCookie(tmpcookie)

	return
}

func (sess *Session) getCookie(r *http.Request, name, secret string) (cookie *Cookie, err error) {
	tmpcookie, err := Helpers.GetCookie(r.Cookies(), name)
	if err != nil {
		return
	}

	value, err := url.QueryUnescape(tmpcookie.Value)
	if err != nil {
		return
	}

	encoding := NewCookieEncoding(name, secret)

	return encoding.Decrypt(value)
}

func (sess *Session) delCookie(w http.ResponseWriter, name string) {
	// force expire client cookie
	tmpcookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Secure:   sess.config.CookieSecure,
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Second),
		MaxAge:   -1,
	}

	http.SetCookie(w, tmpcookie)

	return
}
