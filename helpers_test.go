package session

import (
	"net/http"
	"testing"

	"github.com/golib/assert"
	uuid "github.com/satori/go.uuid"
)

func Test_Helpers_GetCookie(t *testing.T) {
	request, _ := http.NewRequest("HEAD", "http://example.com", nil)

	cookie := &http.Cookie{
		Name:     "getcookie",
		Value:    "cookie",
		Path:     "/",
		Domain:   "example.com",
		HttpOnly: true,
	}
	request.AddCookie(cookie)

	dupcookie := &http.Cookie{
		Name:     "getcookie",
		Value:    "dupcookie",
		Path:     "/",
		Domain:   "example.com",
		HttpOnly: true,
	}
	request.AddCookie(dupcookie)

	// should work
	tmpcookie, err := Helpers.GetCookie(request.Cookies(), "getcookie")
	assert.Nil(t, err)
	assert.Equal(t, tmpcookie.Value, dupcookie.Value)

	// error when not found
	tmpcookie, err = Helpers.GetCookie(request.Cookies(), uuid.NewV4().String())
	assert.EqualError(t, ErrNotFound, err.Error())
	assert.Nil(t, tmpcookie)
}
