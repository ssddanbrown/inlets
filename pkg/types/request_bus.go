package types

import (
	"net/http"
	"sync"
	"time"
)

type RequestBus struct {
	Subscriptions map[string]RequestSubscription
	Mutex         *sync.RWMutex
}

type RequestSubscription struct {
	Data    chan *http.Request
	Created *time.Time
}

func NewRequestBus() *RequestBus {
	return &RequestBus{
		Subscriptions: make(map[string]RequestSubscription),
		Mutex:         &sync.RWMutex{},
	}
}

func (b *RequestBus) SubscriptionList() []string {
	keys := []string{}

	b.Mutex.RLock()
	defer b.Mutex.RUnlock()

	for key := range b.Subscriptions {
		keys = append(keys, key)
	}

	return keys
}

func (b *RequestBus) Send(res *http.Request) {
	var ok bool

	b.Mutex.RLock()
	defer b.Mutex.RUnlock()

	for _, id := range b.SubscriptionList() {
		_, ok = b.Subscriptions[id]

		if !ok {
			continue
		}

		b.Subscriptions[id].Data <- res
	}

}

func (b *RequestBus) Subscribe(id string) *RequestSubscription {
	now := time.Now()
	sub := RequestSubscription{
		Data:    make(chan *http.Request),
		Created: &now,
	}

	b.Mutex.Lock()
	b.Subscriptions[id] = sub
	b.Mutex.Unlock()

	return &sub
}

func (b *RequestBus) Expired(id string, after time.Duration) bool {

	b.Mutex.RLock()

	sub, ok := b.Subscriptions[id]

	b.Mutex.RUnlock()

	if ok && sub.Created != nil {
		return (*sub.Created).Add(after).After(time.Now())
	}

	return false
}

func (b *RequestBus) Unsubscribe(id string) {

	b.Mutex.RLock()

	sub, ok := b.Subscriptions[id]

	b.Mutex.RUnlock()

	if ok {
		close(sub.Data)

		b.Mutex.Lock()
		defer b.Mutex.Unlock()

		delete(b.Subscriptions, id)
	}
}
