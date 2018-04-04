package providers

import (
	"encoding/json"

	"github.com/dolab/session"
	"github.com/go-redis/redis"
)

type Client struct {
	client *redis.Client
}

func New(client *redis.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) New(sid string) (sto session.Storer, err error) {
	sess := NewSessionModel(sid, c.client)

	data, err := json.Marshal(sess.Value)
	if err != nil {
		return
	}

	err = c.client.Set(sess.SID, data, 0).Err()
	if err != nil {
		return
	}

	sto = sess

	return
}

func (c *Client) Restore(sid string) (sto session.Storer, err error) {
	data, err := c.client.Get(sid).Bytes()
	if err != nil {
		return
	}

	var v *session.Value

	err = json.Unmarshal(data, &v)
	if err != nil {
		return
	}

	sto = &SessionModel{
		Client: c.client,
		SID:    sid,
		Value:  v,
	}
	return
}

func (c *Client) Refresh(sid, newsid string) (sto session.Storer, err error) {
	err = c.client.Rename(sid, newsid).Err()
	if err != nil {
		return
	}

	return c.Restore(newsid)
}

func (c *Client) Destroy(sid string) error {
	return c.client.Del(sid).Err()
}
