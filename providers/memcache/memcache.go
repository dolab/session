package providers

import (
	"encoding/json"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dolab/session"
)

type Client struct {
	client *memcache.Client
}

func New(client *memcache.Client) *Client {
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

	item := &memcache.Item{
		Key:   sess.SID,
		Value: data,
	}

	err = c.client.Set(item)
	if err != nil {
		return
	}

	sto = sess

	return
}

func (c *Client) Restore(sid string) (sto session.Storer, err error) {
	item, err := c.client.Get(sid)
	if err != nil {
		return
	}

	var v *session.Value

	err = json.Unmarshal(item.Value, &v)
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
	sto, err = c.Restore(newsid)
	if err != nil {
		return
	}

	err = c.client.Delete(sid)

	return
}

func (c *Client) Destroy(sid string) error {
	return c.client.Delete(sid)
}
