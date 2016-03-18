package cache

import (
	"sync"
	"time"
)

// template type Cache(K,V)
type K int
type V int

var (
	pool sync.Pool
)

type innerElement struct {
	elem         *V
	lastActivity time.Time
}

func newInnerElement(a *V) *innerElement {
	polled_item := pool.Get()
	if polled_item == nil {
		return &innerElement{a, time.Now()}
	}
	ret := polled_item.(*innerElement)
	ret.elem = a
	ret.lastActivity = time.Now()
	return ret
}

func releaseElement(a *innerElement) {
	a.elem = nil
	pool.Put(a)
}

type Cache struct {
	sync.Mutex
	cache         map[K]*innerElement
	removeTimeout time.Duration
}

func NewCache(rt time.Duration) *Cache {
	ret := &Cache{
		cache: make(map[K]*innerElement),
	}
	return ret
}

func (cc *Cache) checkForRemove(k K) {
	for {
		time.Sleep(cc.removeTimeout)
		if cc.tryRemove(k) {
			return
		}
	}
}

func (cc *Cache) tryRemove(k K) bool {
	cc.Lock()
	defer cc.Unlock()

	if inner, ok := cc.cache[k]; ok {
		if time.Now().Sub(inner.lastActivity) >= cc.removeTimeout {
			delete(cc.cache, k)
			releaseElement(inner)
			return true
		}
	}
	return false
}

func (cc *Cache) Set(k K, v *V) {
	cc.Lock()
	defer cc.Unlock()

	if inner, ok := cc.cache[k]; ok {
		inner.elem = v
		inner.lastActivity = time.Now()
	} else {
		cc.cache[k] = newInnerElement(v)
		go cc.checkForRemove(k)
	}
}

func (cc *Cache) Get(k K) *V {
	cc.Lock()
	defer cc.Unlock()

	if inner, ok := cc.cache[k]; ok {
		inner.lastActivity = time.Now()
		return inner.elem
	}
	return nil
}
