package redis

import (
	"encoding/json"
	"time"

	"github.com/dolab/session"
	"github.com/go-redis/redis"
)

type _Session struct{}

var (
	Session *_Session

	db = 0

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // no password set
		Password: "",               // use default DB
		DB:       db,
	})
)

type SessionModel struct {
	SID     string
	Value   *session.Value
	Expires time.Time
}

func NewSessionModel(sid string) *SessionModel {
	return &SessionModel{
		SID:   sid,
		Value: session.NewValue(),
	}
}

func (_ *_Session) New(sid string) (sto session.Storer, err error) {
	sess := NewSessionModel(sid)

	data, err := json.Marshal(sess.Value)
	if err != nil {
		return
	}

	err = client.Set(sess.SID, data, 0).Err()
	if err != nil {
		return
	}

	sto = sess

	return
}

func (_ *_Session) Restore(sid string) (sto session.Storer, err error) {
	data, err := client.Get(sid).Bytes()
	if err != nil {
		return
	}

	var v *session.Value

	err = json.Unmarshal(data, &v)
	if err != nil {
		return
	}

	sto = &SessionModel{
		SID:   sid,
		Value: v,
	}
	return
}

func (_ *_Session) Refresh(sid, newsid string) (sto session.Storer, err error) {

	return
}

func (_ *_Session) Destroy(sid string) error {
	return client.Del(sid).Err()
}

func (session *SessionModel) SessionID() string {
	return session.SID
}

func (session *SessionModel) SetValue(v *session.Value) error {
	session.Value = v

	if v.IsChanged() {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}

		return client.Set(session.SID, data, 0).Err()
	}

	return nil
}

func (session *SessionModel) GetValue() *session.Value {
	return session.Value
}

func (session *SessionModel) Touch() error {
	return nil
}
