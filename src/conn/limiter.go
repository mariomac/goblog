package conn

import (
	"bytes"
	"container/list"
	"fmt"
	"hash/fnv"
	"math"
	"net/http"
	"sync"

	"github.com/mariomac/guara/pkg/rate"

	"time"
)

var clock = time.Now

func ClientRateLimitHandler(inner http.HandlerFunc, maxReqs int, period, clientExpiry time.Duration) http.HandlerFunc {
	retryAfterVal := fmt.Sprint(int(math.Ceil(0.1 * period.Seconds())))
	limiter := NewLimiter(maxReqs, period, clientExpiry)
	return func(rw http.ResponseWriter, req *http.Request) {
		// assumes we are using the http.Server, which sets RemoteAddr to IP:port
		rAddr := []byte(req.RemoteAddr)
		li := bytes.LastIndexByte(rAddr, ':')
		// we don't check li < 0 as long as we are sure the server package adds :port
		if !limiter.Accept(rAddr[:li]) {
			rw.WriteHeader(http.StatusTooManyRequests)
			rw.Header().Add("Retry-After", retryAfterVal)
			// TODO: log rejection
		} else {
			inner(rw, req)
		}
	}
}

// To avoid storing all the possible IPs, we hash it to 65K concurrent IPs
// TODO: make size configurable for larger sites
type IPHash uint16

type Limiter struct {
	maxReqs float64
	period  time.Duration

	cache *TimedLRU[IPHash, *rate.Accepter, bool]
}

func NewLimiter(maxReqs int, period, clientExpiry time.Duration) *Limiter {
	return &Limiter{
		maxReqs: float64(maxReqs),
		period:  period,
		cache:   NewTimedLRU[IPHash, *rate.Accepter, bool](clientExpiry),
	}
}

func (l *Limiter) ClearOldEntries() {
	l.cache.RemoveOldItems()
}

func (l *Limiter) Accept(ip []byte) bool {
	hasher := fnv.New32()
	_, _ = hasher.Write(ip)
	hash32 := hasher.Sum32()
	ipHash := IPHash(hash32>>16 ^ hash32)

	return l.cache.QueryUpdateOrCreate(ipHash, func(r **rate.Accepter) bool {
		return (*r).Accept()
	}, func() (*rate.Accepter, bool) {
		return rate.NewAccepter(l.maxReqs, l.period), true
	})
}

type TimedLRU[K comparable, V, Q any] struct {
	mt      sync.Mutex
	maxTime time.Duration
	ll      *list.List
	cache   map[K]*list.Element
}

type entry[K comparable, V any] struct {
	key      K
	value    V
	lastTime time.Time
}

func NewTimedLRU[K comparable, V, Q any](maxTime time.Duration) *TimedLRU[K, V, Q] {
	return &TimedLRU[K, V, Q]{
		maxTime: maxTime,
		ll:      list.New(),
		cache:   map[K]*list.Element{},
	}
}

func (c *TimedLRU[K, V, Q]) QueryUpdateOrCreate(key K, queryUpdate func(*V) Q, create func() (V, Q)) Q {
	c.mt.Lock()
	defer c.mt.Unlock()
	ee, ok := c.cache[key]
	if ok {
		ee.Value.(*entry[K, V]).lastTime = clock()
		return queryUpdate(&ee.Value.(*entry[K, V]).value)
	}
	cr, qu := create()
	ele := c.ll.PushFront(&entry[K, V]{key: key, value: cr, lastTime: clock()})
	c.cache[key] = ele
	return qu
}

func (c *TimedLRU[K, V, Q]) RemoveOldItems() {
	oldDate := clock().Add(-c.maxTime)
	for c.removeIfOld(oldDate) {
		// hello!
	}
}

func (c *TimedLRU[K, V, Q]) removeIfOld(oldDate time.Time) bool {
	c.mt.Lock()
	defer c.mt.Unlock()

	last := c.ll.Back()
	if last == nil {
		return false
	}

	ve := last.Value.(*entry[K, V])
	if ve.lastTime.After(oldDate) {
		return false
	}

	c.ll.Remove(last)
	delete(c.cache, ve.key)
	return true
}
