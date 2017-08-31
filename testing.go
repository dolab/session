package session

import (
	"log"
	"sync"
	"time"
)

type TestingMemStore struct {
	SID   string
	Data  *Value
	Ctime time.Time
	Atime time.Time
}

func NewTestingMemStore(sid string) *TestingMemStore {
	log.Println("[WARN] It's ONLY used for testing cases.")

	return &TestingMemStore{
		SID:   sid,
		Data:  NewValue(),
		Ctime: time.Now(),
		Atime: time.Now(),
	}
}

func (ms *TestingMemStore) SessionID() string {
	return ms.SID
}

func (ms *TestingMemStore) SetValue(v *Value) error {
	ms.Data = v

	return nil
}

func (ms *TestingMemStore) GetValue() *Value {
	return ms.Data
}

func (ms *TestingMemStore) Touch() error {
	ms.Atime = time.Now()

	return nil
}

type TestingMemProvider struct {
	mux sync.Mutex

	stores map[string]*TestingMemStore
}

func NewTestingMemProvider() *TestingMemProvider {
	log.Println("[WARN] It's ONLY used for testing cases.")

	return &TestingMemProvider{
		stores: map[string]*TestingMemStore{},
	}
}

func (mp *TestingMemProvider) New(sid string) (Storer, error) {
	mp.mux.Lock()
	defer mp.mux.Unlock()

	mp.stores[sid] = NewTestingMemStore(sid)

	return mp.stores[sid], nil
}

func (mp *TestingMemProvider) Restore(sid string) (Storer, error) {
	mp.mux.Lock()
	defer mp.mux.Unlock()

	sto, ok := mp.stores[sid]
	if !ok {
		return nil, ErrNotFound
	}

	return sto, nil
}

func (mp *TestingMemProvider) Refresh(sid, newsid string) (Storer, error) {
	mp.mux.Lock()
	defer mp.mux.Unlock()

	sto, err := mp.Restore(sid)
	if err != nil {
		return nil, err
	}

	newsto := NewTestingMemStore(newsid)
	newsto.SetValue(sto.GetValue())

	mp.stores[newsid] = newsto

	delete(mp.stores, sid)

	return newsto, nil
}

func (mp *TestingMemProvider) Destroy(sid string) error {
	delete(mp.stores, sid)

	return nil
}
