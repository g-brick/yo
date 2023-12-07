# yo 
![example workflow](https://github.com/g-brick/yo/actions/workflows/CI.yml/badge.svg)
[![codecov](https://codecov.io/gh/g-brick/yo/graph/badge.svg?token=8AGHULWWDJ)](https://codecov.io/gh/g-brick/yo)
[![Go Report Card](https://goreportcard.com/badge/github.com/g-brick/yo)](https://goreportcard.com/report/github.com/g-brick/yo)

**Yo** is a wrapper sync.Group library, which is very simple to use compared to sync.Group, 
and it is very lightweight and has almost no external third-party dependencies. 
You don't have to worry about when you add or remove the number of goroutines, 
using it will make your concurrent programming or asynchronous tasks more concise and elegant.
## Getting started
Simply add the following import in your code file, 
and then `go [build|run|test]` will automatically fetch the necessary dependencies.
```
import "github.com/g-brick/yo"
```
Otherwise, run the following Go command to install the `yo` package:

```sh
$ go get -u github.com/g-brick/yo
```
## How to use
### General usage
First you need to import Yo package for using Yo, one example :
```go
package main
import (
	"context"
	"fmt"
	"sync/atomic"
	"github.com/g-brick/yo"
)
func main() {
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
```
### Usage with a limited number of goroutines
You also can use GOMAXPROCS in yo like this to limit the number of goroutines.
```go
package main
import (
	"context"
	"fmt"
	"sync/atomic"
	"github.com/g-brick/yo"
)
func main() {
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
```
### Usage with a timeout control
If you want to timeout control multiple tasks without having to wait for them all to succeed, 
you can use this method to avoid delaying your own application processing time 
when requesting multiple external third-party applications in parallel.

```go
package main
import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
	"github.com/g-brick/yo"
)
// Used by timeout control.
func main() {
	var (
		urls      = []string{"https://bing.com", "https://github.com", "https://google.com", "https://baidu.com", "https://stackoverflow.com"}
		okReq     int32
		deadline  = time.Millisecond * 1200 // 1.2s
		c, cancel = context.WithTimeout(context.Background(), deadline)
	)
	defer cancel()
	y := yo.WithCancel(c)
	for _, url := range urls { // 5 requests concurrent in bing.com would be canceled if some requests were timeout.
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
```
### Usage in fan-out mode
In fan-out mode, you can run a large number of asynchronous tasks in the background. 
You can set the number of workers and the length of the task queue. 
It avoids the frenzied creation and destruction of goroutines.
It should be noted that when the buffer is full, new processing tasks will not be processed, 
so you need to consider the ratio of worker to buffer.
```go
package main
import (
	"context"
	"fmt"
	"time"
	"github.com/g-brick/yo"
)
func main() {
	taskDealer := yo.NewFanout("TaskDealer", yo.Worker(50), yo.Buffer(1024)) // Set up a global dealer with 50 goroutines to handle a 1024-length queue in the background.
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
```


