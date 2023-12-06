package examples

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/g-brick/yo"
)

func _normalUsage() {
	var (
		count int32
		c     = context.Background()
	)
	y := yo.WithContext(c)
	for i := 0; i < 100; i++ {
		y.Go(func(ctx context.Context) (err error) {
			atomic.AddInt32(&count, 1)
			return
		})
	}
	if err := y.Wait(); err != nil { // Wait the completion of every goroutine.
		fmt.Printf("Err is %v", err)
		return
	}
	fmt.Printf("The final count is %d", count) // The count must be 100.
}

// Used by timeout control.
func _withCancelControl() {
	var (
		url       = "https://bing.com"
		okReq     int32
		deadline  = time.Millisecond * 1200 // 1.2s
		c, cancel = context.WithTimeout(context.Background(), deadline)
	)
	defer cancel()
	y := yo.WithCancel(c)
	for i := 0; i < 5; i++ { // 5 requests concurrent in bing.com would be canceled if some requests were timeout.
		y.Go(func(ctx context.Context) (err error) {
			if _, err = http.Get(url); err == nil {
				atomic.AddInt32(&okReq, 1)
			}
			return
		})
	}
	if err := y.Wait(); err != nil { // Wait the completion of every goroutine.
		fmt.Printf("Err is %v", err)
	}
	fmt.Printf("5 requests, and %d request(s) succeeded", okReq)
}

// Used in a normal usage of the errGroup with yo by a limited version.
func _WithLimitedGoroutines() {
	var (
		count int32
		c     = context.Background()
	)
	y := yo.WithContext(c)
	y.GOMAXPROCS(5) // Limit the nums of goroutine here.
	for i := 0; i < 100; i++ {
		y.Go(func(ctx context.Context) (err error) {
			atomic.AddInt32(&count, 1)
			return
		})
	}
	if err := y.Wait(); err != nil { // Wait the completion of every goroutine.
		fmt.Printf("Err is %v", err)
		return
	}
	fmt.Printf("The final count is %d", count) // The count must be 100.
}
