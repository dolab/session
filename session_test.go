package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golib/assert"
	uuid "github.com/satori/go.uuid"
)

var (
	config = &Config{
		CookieName:   "testingcookie",
		CookieExpire: 3600, // 1h
		CookieSecure: false,
		CookieSecret: uuid.NewV4().String(),
	}

	provider = NewTestingMemProvider()
)

func Test_Session_New(t *testing.T) {
	assertion := assert.New(t)
	session := New(provider, config)

	request, _ := http.NewRequest("HEAD", "http://example.com", nil)
	response := httptest.NewRecorder()
	assertion.Empty(response.HeaderMap)

	// should work
	sto, err := session.New(response, request)
	assertion.Nil(err)
	assertion.NotEmpty(response.HeaderMap["Set-Cookie"])

	tmpsto, err := provider.Restore(sto.SessionID())
	assertion.Nil(err)
	assertion.Equal(sto, tmpsto)
}

func Test_Session_Start(t *testing.T) {
	assertion := assert.New(t)

	request, _ := http.NewRequest("HEAD", "http://example.com", nil)
	response := httptest.NewRecorder()

	// register cookie with provider
	tmpsto, _ := provider.New(uuid.NewV4().String())
	tmpval := tmpsto.GetValue()
	tmpval.Set("key", "value")

	// inject cookie
	cookie := NewCookie(tmpsto.SessionID())
	data, _ := NewCookieEncoding(config.CookieName, config.CookieSecret).Encrypt(cookie)

	request.AddCookie(&http.Cookie{
		Name:   config.CookieName,
		Value:  data,
		Path:   "/",
		Domain: "example.com",
		MaxAge: config.CookieExpire,
	})

	// should work
	sess := New(provider, config)

	sto, err := sess.Start(response, request)
	assertion.Nil(err)
	assertion.Empty(response.HeaderMap)
	assertion.Equal(tmpsto.SessionID(), sto.SessionID())

	val := sto.GetValue()
	assertion.Equal(tmpval.String("key"), val.String("key"))
}

func Test_Session_Refresh(t *testing.T) {
	assertion := assert.New(t)

	request, _ := http.NewRequest("HEAD", "http://example.com", nil)
	response := httptest.NewRecorder()

	// register cookie with provider
	tmpsto, _ := provider.New(uuid.NewV4().String())
	tmpval := tmpsto.GetValue()
	tmpval.Set("key", "value")

	// inject cookie
	cookie := NewCookie(tmpsto.SessionID())
	data, _ := NewCookieEncoding(config.CookieName, config.CookieSecret).Encrypt(cookie)

	request.AddCookie(&http.Cookie{
		Name:   config.CookieName,
		Value:  data,
		Path:   "/",
		Domain: "example.com",
		MaxAge: config.CookieExpire,
	})

	// should work
	sess := New(provider, config)

	sto, err := sess.Refresh(response, request)
	assertion.Nil(err)
	assertion.NotEmpty(response.HeaderMap["Set-Cookie"])
	assertion.NotEqual(tmpsto.SessionID(), sto.SessionID())

	val := sto.GetValue()
	assertion.Equal(tmpval.String("key"), val.String("key"))

	oldsto, err := provider.Restore(tmpsto.SessionID())
	assertion.EqualError(ErrNotFound, err.Error())
	assertion.Nil(oldsto)

	newsto, err := provider.Restore(sto.SessionID())
	assertion.Nil(err)
	assertion.Equal(sto, newsto)
}

func Test_Session_Destroy(t *testing.T) {
	assertion := assert.New(t)

	request, _ := http.NewRequest("HEAD", "http://example.com", nil)
	response := httptest.NewRecorder()

	// register cookie with provider
	tmpsto, _ := provider.New(uuid.NewV4().String())
	tmpval := tmpsto.GetValue()
	tmpval.Set("key", "value")

	// inject cookie
	cookie := NewCookie(tmpsto.SessionID())
	data, _ := NewCookieEncoding(config.CookieName, config.CookieSecret).Encrypt(cookie)

	request.AddCookie(&http.Cookie{
		Name:   config.CookieName,
		Value:  data,
		Path:   "/",
		Domain: "example.com",
		MaxAge: config.CookieExpire,
	})

	// should work
	sess := New(provider, config)

	err := sess.Destroy(response, request)
	assertion.Nil(err)
	assertion.NotEmpty(response.HeaderMap["Set-Cookie"])

	oldsto, err := provider.Restore(tmpsto.SessionID())
	assertion.EqualError(ErrNotFound, err.Error())
	assertion.Nil(oldsto)
}
