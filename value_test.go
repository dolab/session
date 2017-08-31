package session

import (
	"testing"
	"time"

	"github.com/golib/assert"
	uuid "github.com/satori/go.uuid"
)

func Test_Value(t *testing.T) {
	value := NewValue()
	assert.Implements(t, (*Valuer)(nil), value)
	assert.False(t, value.IsChanged())
}

func Test_Value_Add(t *testing.T) {
	value := NewValue()

	key := uuid.NewV4().String()
	val := time.Now()

	// should work
	assert.False(t, value.Has(key))

	err := value.Add(key, val)
	assert.Nil(t, err)
	assert.True(t, value.Has(key))
	assert.True(t, value.IsChanged())

	// error when duplicated
	err = value.Add(key, val)
	assert.EqualError(t, ErrDuplicateKey, err.Error())
}

func Test_Value_Set(t *testing.T) {
	value := NewValue()

	key := uuid.NewV4().String()
	val := time.Now()

	// should work
	assert.False(t, value.Has(key))

	err := value.Set(key, val)
	assert.Nil(t, err)
	assert.True(t, value.Has(key))
	assert.True(t, value.IsChanged())

	// should work even duplicated
	err = value.Set(key, val)
	assert.Nil(t, err)
}

func Test_Value_Get(t *testing.T) {
	value := NewValue()
	key := uuid.NewV4().String()
	val := time.Now()

	value.Add(key, val)

	// should work for Get
	tmp, err := value.Get(key)
	assert.Nil(t, err)
	assert.NotEmpty(t, tmp)
}

func Test_Value_Unmarshal(t *testing.T) {
	value := NewValue()

	key := uuid.NewV4().String()
	val := time.Now()

	value.Add(key, val)

	// should work for Unmarshal
	var lav time.Time

	err := value.Unmarshal(key, &lav)
	assert.Nil(t, err)
	assert.Equal(t, val.Unix(), lav.Unix())
}

func Test_Value_Del(t *testing.T) {
	value := NewValue()
	key := uuid.NewV4().String()
	val := time.Now()

	value.Add(key, val)

	// should work for Del
	err := value.Del(key)
	assert.Nil(t, err)
	assert.False(t, value.Has(key))

	// should work even without a key
	err = value.Del(uuid.NewV4().String())
	assert.Nil(t, err)
}
