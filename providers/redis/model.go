package providers

import (
	"encoding/json"
	"time"

	"github.com/dolab/session"
	"github.com/go-redis/redis"
)

type SessionModel struct {
	Client  *redis.Client
	SID     string
	Value   *session.Value
	Expires time.Time
}

func NewSessionModel(sid string, client *redis.Client) *SessionModel {
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

		return this.Client.Set(this.SID, data, 0).Err()
	}

	return nil
}

func (this *SessionModel) GetValue() *session.Value {
	return this.Value
}

func (this *SessionModel) Touch() error {
	return nil
}
