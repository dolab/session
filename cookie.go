package session

import (
	"bytes"
	"encoding/gob"
	"strings"
	"time"
)

type Cookie struct {
	Value string    `gob:"v"`
	Ctime time.Time `gob:"ct"`
}

func NewCookie(value string) *Cookie {
	return &Cookie{
		Value: value,
		Ctime: time.Now(),
	}
}

func (cookie *Cookie) IsExpiredBy(seconds time.Duration) bool {
	return time.Now().After(cookie.Ctime.Add(seconds))
}

type CookieEncoding struct {
	encoding *Encoding
}

func NewCookieEncoding(cookieName, cookieSecret string) *CookieEncoding {
	return &CookieEncoding{
		encoding: NewEncoding(cookieName, cookieSecret),
	}
}

func (ce *CookieEncoding) Encrypt(cookie *Cookie) (data string, err error) {
	buf := bytes.NewBuffer(nil)

	err = gob.NewEncoder(buf).Encode(cookie)
	if err != nil {
		return
	}

	data = ce.encoding.Encrypt(buf.String())
	return
}

func (ce *CookieEncoding) Decrypt(data string) (cookie *Cookie, err error) {
	s, err := ce.encoding.Decrypt(data)
	if err != nil {
		return
	}

	err = gob.NewDecoder(strings.NewReader(s)).Decode(&cookie)
	return
}
