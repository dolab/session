package session

import (
	"encoding/json"
	"sync"
)

// Value defines session data container
// TODO: refactor to use go concurrent map?
type Value struct {
	mux     sync.RWMutex
	Data    map[string]json.RawMessage
	changed bool
}

func NewValue() *Value {
	return &Value{
		Data:    map[string]json.RawMessage{},
		changed: false,
	}
}

func (v *Value) IsChanged() bool {
	return v.changed
}

func (v *Value) Has(key string) bool {
	v.mux.RLock()
	_, ok := v.Data[key]
	v.mux.RUnlock()

	return ok
}

func (v *Value) Unmarshal(key string, value interface{}) (err error) {
	data, err := v.Get(key)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &value)
	return
}

func (v *Value) Add(key string, value interface{}) error {
	if v.Has(key) {
		return ErrDuplicateKey
	}

	return v.Set(key, value)
}

func (v *Value) Set(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	v.mux.Lock()
	v.Data[key] = json.RawMessage(b)
	v.changed = true
	v.mux.Unlock()

	return nil
}

func (v *Value) Get(key string) (data json.RawMessage, err error) {
	v.mux.RLock()
	data, ok := v.Data[key]
	if !ok {
		err = ErrNotFound
	}
	v.mux.RUnlock()

	return
}

func (v *Value) Del(key string) error {
	delete(v.Data, key)
	v.changed = true

	return nil
}

func (v *Value) Clean() {
	v.Data = map[string]json.RawMessage{}
	v.changed = true

	return
}

func (v *Value) Bool(key string) (ok bool, err error) {
	err = v.Unmarshal(key, &ok)
	if err == nil {
		return
	}

	// fallback to weakness
	switch v.String(key) {
	case "1", "t", "true", "on", "y", "yes":
		ok = true
		err = nil
		return
	}

	return
}

func (v *Value) Int(key string) (i int, err error) {
	err = v.Unmarshal(key, &i)
	return
}

func (v *Value) Int8(key string) (i8 int8, err error) {
	err = v.Unmarshal(key, &i8)
	return
}

func (v *Value) Int16(key string) (i16 int16, err error) {
	err = v.Unmarshal(key, &i16)
	return
}

func (v *Value) Int32(key string) (i32 int32, err error) {
	err = v.Unmarshal(key, &i32)
	return
}

func (v *Value) Int64(key string) (i64 int64, err error) {
	err = v.Unmarshal(key, &i64)
	return
}

func (v *Value) Uint(key string) (u uint, err error) {
	err = v.Unmarshal(key, &u)
	return
}

func (v *Value) Uint8(key string) (u8 uint8, err error) {
	err = v.Unmarshal(key, &u8)
	return
}

func (v *Value) Uint16(key string) (u16 uint16, err error) {
	err = v.Unmarshal(key, &u16)
	return
}

func (v *Value) Uint32(key string) (u32 uint32, err error) {
	err = v.Unmarshal(key, &u32)
	return
}

func (v *Value) Uint64(key string) (u64 uint64, err error) {
	err = v.Unmarshal(key, &u64)
	return
}

func (v *Value) Float32(key string) (f32 float32, err error) {
	err = v.Unmarshal(key, &f32)
	return
}

func (v *Value) Float64(key string) (f64 float64, err error) {
	err = v.Unmarshal(key, &f64)
	return
}

func (v *Value) String(key string) string {
	data, err := v.Get(key)
	if err != nil {
		return ""
	}

	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(b)
}
