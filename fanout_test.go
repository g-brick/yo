package yo

import (
	"context"
	"testing"
	"time"
)

func TestFanoutDo(t *testing.T) {
	ca := NewFanout("back_proc", Worker(1), Buffer(1024))
	var run bool
	ca.Do(context.Background(), func(c context.Context) {
		run = true
		panic("error")
	})
	time.Sleep(time.Millisecond * 50)
	t.Log("not panic in main")
	if !run {
		t.Fatal("expect run be true")
	}
}

func TestFanoutClose(t *testing.T) {
	ca := NewFanout("back_proc", Worker(1), Buffer(1024))
	ca.Close()
	err := ca.Do(context.Background(), func(c context.Context) {})
	if err == nil {
		t.Fatal("expect get err")
	}
}
