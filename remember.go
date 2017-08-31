package session

import (
	"net/http"
)

func (sess *Session) NewRemember(w http.ResponseWriter, r *http.Request, psk string) (cookie *Cookie, err error) {
	return sess.setCookie(w, r, sess.config.RememberName, psk, sess.config.RememberSecret, sess.config.RememberExpire)
}

func (sess *Session) StartRemember(r *http.Request) (psk string, err error) {
	cookie, err := sess.getCookie(r, sess.config.RememberName, sess.config.RememberSecret)
	if err != nil {
		return
	}

	if cookie.IsExpiredBy(sess.config.RememberExpireSeconds()) {
		err = ErrCookieExpired
		return
	}

	psk = cookie.Value

	return
}

func (sess *Session) DestroyRemember(w http.ResponseWriter, r *http.Request) error {
	// force expire client cookie
	sess.delCookie(w, sess.config.RememberName)

	return nil
}
