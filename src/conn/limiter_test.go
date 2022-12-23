package conn

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewLimiter(t *testing.T) {
	l := NewLimiter(60, 100 * time.Millisecond, time.Minute)

	ip1 := net.IPv4(1,2,3,4)
	ip2 := net.IPv4(4,3,2,1)

	start := time.Now()

	// should accept all the queries as long as there are less than 60/100ms
	for time.Now().Sub(start) < 30 * time.Millisecond {
		assert.True(t, l.Accept(ip1))
		time.Sleep(time.Millisecond)
	}
	// should end up rejecting queries
	last := true
	for time.Now().Sub(start) < 90 * time.Millisecond {
		last = last && l.Accept(ip1)
	}
	assert.False(t, last)

	time.Sleep(30 * time.Millisecond)
	l.ClearOldEntries()

	// after some time, there is more room to get enough queries
	start = time.Now()
	for i := 0 ; i < 10 ; i++ {
		assert.True(t, l.Accept(ip1))
		assert.True(t, l.Accept(ip2))
		time.Sleep(time.Millisecond)
	}
}

func TestTimedLRU(t *testing.T) {
	now := time.Now()
	clock = func() time.Time {
		return now
	}
	tlru := NewTimedLRU[string, string, int](time.Second)

	qu := func(i *string) int {
		*i = *i + "-"
		return len(*i)
	}
	cr := func () (string, int) {
		return "", 0
	}
	assert.Equal(t, 0, tlru.QueryUpdateOrCreate("hi", qu, cr))
	assert.Equal(t, 0, tlru.QueryUpdateOrCreate("ho", qu, cr))
	now = now.Add(800 * time.Millisecond)
	assert.Equal(t, 0, tlru.QueryUpdateOrCreate("foo", qu, cr))
	assert.Equal(t, 1, tlru.QueryUpdateOrCreate("ho", qu, cr))
	now = now.Add(800 * time.Millisecond)
	tlru.RemoveOldItems()
	// "hi" has been deleted then it returns a fresh one
	assert.Equal(t, 0, tlru.QueryUpdateOrCreate("hi", qu, cr))
	assert.Equal(t, 2, tlru.QueryUpdateOrCreate("ho", qu, cr))
	assert.Equal(t, 1, tlru.QueryUpdateOrCreate("foo", qu, cr))
	now = now.Add(10000 * time.Millisecond)
	tlru.RemoveOldItems()
	// all entries have been deleted then it returns fresh instances
	assert.Equal(t, 0, tlru.QueryUpdateOrCreate("hi", qu, cr))
	assert.Equal(t, 0, tlru.QueryUpdateOrCreate("ho", qu, cr))
	assert.Equal(t, 0, tlru.QueryUpdateOrCreate("foo", qu, cr))
}

func TestTimedLRU_ConcurrencyRaces(t *testing.T) {
	timeBias := int64(0)
	clock = func() time.Time {
		return time.Now().Add(time.Duration(atomic.LoadInt64(&timeBias)))
	}
	tlru := NewTimedLRU[string, int64, bool](time.Second)
	wg := sync.WaitGroup{}
	wg.Add(4)
	for i:= 0 ; i < 4 ; i++ {
		thread := i
		go func() {
			start := time.Now()
			defer wg.Done()
			for time.Now().Sub(start) < 100 * time.Millisecond {
				clk := clock()
				key := fmt.Sprint(clk.UnixMilli())
				tlru.QueryUpdateOrCreate(key, func(v *int64) bool {
					atomic.AddInt64(v, 1)
					return true
				}, func() (int64, bool) {
					return 0, true
				})
				tlru.RemoveOldItems()
				runtime.Gosched()
			}
			if thread == 0 {
				atomic.StoreInt64(&timeBias, int64(950 * time.Millisecond))
				tlru.RemoveOldItems()
			} else {
				for time.Now().Sub(start) < 200 * time.Millisecond {
					clk := clock()
					key := fmt.Sprint(clk.UnixMilli())
					tlru.QueryUpdateOrCreate(key, func(v *int64) bool {
						atomic.AddInt64(v, 1)
						return true
					}, func() (int64, bool) {
						return 0, true
					})
					tlru.RemoveOldItems()
					runtime.Gosched()
				}
			}
		}()
	}
	wg.Wait()
	fmt.Println("hoe!")
}