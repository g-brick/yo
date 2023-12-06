package yo

import (
	"context"
	"testing"
	"time"
)

func TestFanoutDo(t *testing.T) {
	bg := NewFanout("background execution", Worker(1), Buffer(1024))
	var ok bool
	bg.Do(context.Background(), func(c context.Context) {
		ok = true
		panic("error")
	})
	time.Sleep(time.Millisecond * 50)
	t.Log("not panic in main")
	if !ok {
		t.Fatal("expect ok == true")
	}
}

func TestFanoutClose(t *testing.T) {
	bg := NewFanout("background execution", Worker(1), Buffer(1024))
	_ = bg.Close()
	err := bg.Do(context.Background(), func(c context.Context) {})
	if err == nil {
		t.Fatal("expect a err")
	}
}
