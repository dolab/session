package session

import "encoding/json"

// Provider defines session store provider apis
type Provider interface {
	New(sid string) (sto Storer, err error)
	Restore(sid string) (sto Storer, err error)
	Refresh(sid, newsid string) (Storer, error)
	Destroy(sid string) error
}

// Storer defines session store apis
type Storer interface {
	SessionID() string       // return current session ID
	SetValue(v *Value) error // set session data
	GetValue() *Value        // get sesstion data
	Touch() error            // sync session expire time to the provider
}

type Valuer interface {
	IsChanged() bool                               // whether value has updated?
	Has(key string) bool                           // is key exist?
	Unmarshal(key string, value interface{}) error // unmarshal a value by key with value type
	Add(key string, value interface{}) error       // add a value with the key, returns an error if existed
	Set(key string, value interface{}) error       // set a value with the key
	Get(key string) (json.RawMessage, error)       // get a value by key
	Del(key string) error                          // delete a value by key
	Clean()                                        // clean all values
}
