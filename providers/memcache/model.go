package providers

import (
	"encoding/json"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dolab/session"
)

type SessionModel struct {
	Client  *memcache.Client
	SID     string
	Value   *session.Value
	Expires time.Time
}

func NewSessionModel(sid string, client *memcache.Client) *SessionModel {
	return &SessionModel{
		Client: client,
		SID:    sid,
		Value:  session.NewValue(),
	}
}

func (this *SessionModel) SessionID() string {
	return this.SID
}

func (this *SessionModel) SetValue(v *session.Value) error {
	this.Value = v

	if v.IsChanged() {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}

		item := &memcache.Item{
			Key:   this.SID,
			Value: data,
		}

		return this.Client.Set(item)
	}

	return nil
}

func (this *SessionModel) GetValue() *session.Value {
	return this.Value
}

func (this *SessionModel) Touch() error {
	return nil
}
