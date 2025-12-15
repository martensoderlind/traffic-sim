package events

import (
	"reflect"
	"sync"
)

type Handler func(payload any)

type Dispatcher struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{handlers: make(map[string][]Handler)}
}

func (d *Dispatcher) Subscribe(name string, h Handler) func() {
	d.mu.Lock()
	d.handlers[name] = append(d.handlers[name], h)
	d.mu.Unlock()

	return func() {
		d.mu.Lock()
		hs := d.handlers[name]
		for i := range hs {
			if reflect.ValueOf(hs[i]).Pointer() == reflect.ValueOf(h).Pointer() {
				d.handlers[name] = append(hs[:i], hs[i+1:]...)
				break
			}
		}
		d.mu.Unlock()
	}
}

func (d *Dispatcher) Emit(name string, payload any) {
	d.mu.RLock()
	hs := append([]Handler(nil), d.handlers[name]...)
	d.mu.RUnlock()

	for _, h := range hs {
		h(payload)
	}
}
