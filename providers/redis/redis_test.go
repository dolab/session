package providers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dolab/session"
	"github.com/go-redis/redis"
	"github.com/golib/assert"
	uuid "github.com/satori/go.uuid"
)

var (
	config = &session.Config{
		CookieName:   "testingcookie",
		CookieExpire: 3600, // 1h
		CookieSecure: false,
		CookieSecret: uuid.Must(uuid.NewV4()).String(),
	}
)

func Test_Session_New(t *testing.T) {
	assertion := assert.New(t)

	client := New(redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}))

	sess := session.New(client, config)

	request, _ := http.NewRequest("HEAD", "http://example.com", nil)
	response := httptest.NewRecorder()
	assertion.Empty(response.HeaderMap)

	// should work
	sto, err := sess.New(response, request)
	assertion.Nil(err)
	assertion.NotEmpty(response.HeaderMap["Set-Cookie"])

	tmpsto, err := client.Restore(sto.SessionID())
	assertion.Nil(err)
	assertion.Equal(sto, tmpsto)

	value := sto.GetValue()
	err = value.Set("current_user", "tmpUser")
	assertion.Nil(err)

	err = sto.SetValue(value)
	assertion.Nil(err)

	var user string
	err = sto.GetValue().Unmarshal("current_user", &user)
	assertion.Nil(err)

	assertion.Equal("tmpUser", user)
}
