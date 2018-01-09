package session

import (
	"testing"

	"github.com/golib/assert"
	"github.com/satori/go.uuid"
)

func Test_Encoding(t *testing.T) {
	block := uuid.Must(uuid.NewV4()).String()
	key := uuid.Must(uuid.NewV4()).String()
	value := uuid.Must(uuid.NewV4()).String()

	enc := NewEncoding(block, key)

	// encrypt
	encstr := enc.Encrypt(value)
	assert.NotEmpty(t, encstr)

	// decrypt
	decstr, err := enc.Decrypt(encstr)
	assert.Nil(t, err)
	assert.Equal(t, value, decstr)
}
