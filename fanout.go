package yo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var (
	// ErrFull chan full.
	ErrFull = errors.New("task busy, Fanout: chan is overflow")
)

type options struct {
	worker int
	buffer int
}

// Option Fanout option
type Option func(*options)

// Worker specifies the worker of Fanout
func Worker(n int) Option {
	if n <= 0 {
		panic("Fanout: worker should > 0")
	}
	return func(o *options) {
		o.worker = n
	}
}

// Buffer specifies the buffer of Fanout
func Buffer(n int) Option {
	if n <= 0 {
		panic("Fanout: buffer should > 0")
	}
	return func(o *options) {
		o.buffer = n
	}
}

type item struct {
	f   func(c context.Context)
	ctx context.Context
}

// Fanout async consume data from chan.
type Fanout struct {
	name    string
	ch      chan item
	options *options
	waiter  sync.WaitGroup

	ctx    context.Context
	cancel func()
}

// NewFanout That is new a Fanout struct.
func NewFanout(name string, opts ...Option) *Fanout {
	if name == "" {
		name = "anonymous"
	}
	o := &options{
		worker: 1,
		buffer: 1024,
	}
	for _, op := range opts {
		op(o)
	}
	c := &Fanout{
		ch:      make(chan item, o.buffer),
		name:    name,
		options: o,
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.waiter.Add(o.worker)
	for i := 0; i < o.worker; i++ {
		go c.proc()
	}
	return c
}

func (c *Fanout) proc() {
	defer c.waiter.Done()
	for {
		select {
		case t := <-c.ch:
			wrapFunc(t.f)(t.ctx)
		// TODO Maybe add some monitor metrics in the future here
		case <-c.ctx.Done():
			return
		}
	}
}

func wrapFunc(f func(c context.Context)) (res func(context.Context)) {
	res = func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64*1024)
				buf = buf[:runtime.Stack(buf, false)]
				log.Printf(
					"[wrapFunc] panic in Fanout proc, %v, %v, %v",
					zap.String("err", fmt.Sprint(r)),
					zap.String("fn", getFuncName(f)),
					zap.String("stack", string(buf)),
				)
			}
		}()
		f(ctx)
	}
	return
}

// Do save a callback func.
func (c *Fanout) Do(ctx context.Context, f func(context.Context)) (err error) {
	if f == nil || c.ctx.Err() != nil {
		return c.ctx.Err()
	}
	select {
	case c.ch <- item{f: f, ctx: ctx}:
	default:
		err = ErrFull
	}
	return
}

// Close close Fanout
func (c *Fanout) Close() error {
	if err := c.ctx.Err(); err != nil {
		return err
	}
	c.cancel()
	c.waiter.Wait()
	return nil
}

func getFuncName(i interface{}, seps ...rune) (funcName string) {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})
	size := len(fields)
	if size == 0 {
		return
	}
	funcName = fields[size-1]
	return
}
