package events

import (
	"testing"
	"time"
)

func TestSubscribeEmitUnsubscribe(t *testing.T) {
    d := NewDispatcher()
    ch := make(chan any, 1)

    h := func(p any) { ch <- p }
    unsub := d.Subscribe("test.event", h)
    d.Emit("test.event", "payload")

    select {
    case v := <-ch:
        if v.(string) != "payload" {
            t.Fatalf("expected payload, got %v", v)
        }
    case <-time.After(1 * time.Second):
        t.Fatalf("handler not called")
    }

    unsub()
    d.Emit("test.event", "x")

    select {
    case <-ch:
        t.Fatalf("handler should have been unsubscribed")
    case <-time.After(100 * time.Millisecond):
        // ok
    }
}
