package examples

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/g-brick/yo"
)

func _normalUsage() {
	var (
		count int32
		l     sync.RWMutex
		c     = context.Background()
	)
	y := yo.WithContext(c)
	for i := 0; i < 100; i++ {
		y.Go(func(ctx context.Context) (err error) {
			l.Lock()
			defer l.Unlock()
			count++
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
		urls      = []string{"https://bing.com", "https://github.com", "https://google.com", "https://baidu.com", "https://stackoverflow.com"}
		okReq     int32
		l         sync.RWMutex
		deadline  = time.Millisecond * 1200 // 1.2s
		c, cancel = context.WithTimeout(context.Background(), deadline)
	)
	defer cancel()
	y := yo.WithCancel(c)
	for _, u := range urls { // 5 requests concurrent in bing.com would be canceled if some requests were timeout.
		url := u
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
		fmt.Printf("Err is %v", err)
	}
	fmt.Printf("5 requests, and %d request(s) succeeded", okReq)
}

// Used in a normal usage of the errGroup with yo by a limited version.
func _WithLimitedGoroutines() {
	var (
		count int32
		l     sync.RWMutex
		c     = context.Background()
	)
	y := yo.WithContext(c)
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
		fmt.Printf("Err is %v", err)
		return
	}
	fmt.Printf("The final count is %d", count) // The count must be 100.
}

// Fanout mode usage：makes many tasks run in the background with a controllable numbers of goroutines.
func _withFanoutMode() {
	taskDealer := yo.NewFanout("TaskDealer", yo.Worker(50), yo.Buffer(1024))
	for i := 0; i < 150; i++ {
		err := taskDealer.Do(context.Background(), func(ctx context.Context) {
			// Do something heavy task here asynchronously
			time.Now()
		})
		if err != nil {
			fmt.Printf("Err is %v", err)
		}
	}
}
