package session

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

var (
	Helpers *_Helper
)

type _Helper struct{}

func (_ *_Helper) NewSessionID() (sid string) {
	h := md5.New()
	h.Write(uuid.NewV4().Bytes())

	sid = hex.EncodeToString(h.Sum(nil))
	return
}

// get last same name cookie from cookies
// http://play.golang.org/p/LDfjMnJnhI
func (_ *_Helper) GetCookie(cookies []*http.Cookie, name string) (cookie *http.Cookie, err error) {
	for i := len(cookies) - 1; i >= 0; i-- {
		if cookies[i].Name == name {
			cookie = cookies[i]
			break
		}
	}

	if cookie == nil {
		err = ErrNotFound
	}

	return
}
