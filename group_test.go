package yo

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"
)

// We test it in a normal usage of the errGroup in yo.
func TestNormalUsage(t *testing.T) {
	var (
		count int32
		l     sync.RWMutex
		c     = context.Background()
	)
	y := WithContext(c)
	for i := 0; i < 100; i++ {
		y.Go(func(ctx context.Context) (err error) {
			l.Lock()
			defer l.Unlock()
			count++
			return
		})
	}
	if err := y.Wait(); err != nil { // Wait the completion of every goroutine.
		t.Errorf("Err is %v", err)
		return
	}
	t.Logf("The final count is %d", count) // The count must be 100.
}

// We test it by timeout control.
func TestWithCancel(t *testing.T) {
	var (
		urls      = []string{"https://bing.com", "https://github.com", "https://google.com", "https://baidu.com", "https://stackoverflow.com"}
		okReq     int32
		l         sync.RWMutex
		deadline  = time.Millisecond * 1200 // 1.2s
		c, cancel = context.WithTimeout(context.Background(), deadline)
	)
	defer cancel()
	y := WithCancel(c)
	for _, url := range urls { // 5 requests concurrent in bing.com would be canceled if some requests were timeout.
		y.Go(func(ctx context.Context) (err error) {
			if _, err = http.Get(url); err == nil {
				l.Lock()
				defer l.Unlock()
				okReq++
			}
			return
		})
	}
	if err := y.Wait(); err != nil { // Wait the completion of every goroutine.
		t.Logf("Err is %v", err)
	}
	t.Logf("5 requests, and %d request(s) succeeded", okReq)
}

// We test it in a normal usage of the errGroup with yo by a limited version.
func TestWithLimitedGoroutines(t *testing.T) {
	var (
		count int32
		l     sync.RWMutex
		c     = context.Background()
	)
	y := WithContext(c)
	y.GOMAXPROCS(5) // Limit the nums of goroutine here.
	for i := 0; i < 100; i++ {
		y.Go(func(ctx context.Context) (err error) {
			l.Lock()
			defer l.Unlock()
			count++
			return
		})
	}
	if err := y.Wait(); err != nil { // Wait the completion of every goroutine.
		t.Errorf("Err is %v", err)
		return
	}
	t.Logf("The final count is %d", count) // The count must be 100.
}
