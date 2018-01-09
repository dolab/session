package session

import (
	"testing"
	"time"

	"github.com/golib/assert"
	uuid "github.com/satori/go.uuid"
)

func Test_Cookie(t *testing.T) {
	value := uuid.Must(uuid.NewV4()).String()

	cookie := NewCookie(value)
	assert.Equal(t, value, cookie.Value)
	assert.NotZero(t, cookie.Ctime)
	assert.False(t, cookie.IsExpiredBy(1*time.Second))
}
