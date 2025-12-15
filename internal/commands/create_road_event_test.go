package commands

import (
	"testing"
	"time"
	"traffic-sim/internal/events"
	"traffic-sim/internal/road"
	"traffic-sim/internal/world"
)

func TestCreateRoadEmitsEvent(t *testing.T) {
	w := world.New()
	ex := NewCommandExecutor(w)

	ch := make(chan any, 1)
	w.Events.Subscribe(events.EventRoadCreated, func(p any) { ch <- p })

	from := &road.Node{ID: "n1", X: 0, Y: 0}
	to := &road.Node{ID: "n2", X: 100, Y: 0}
	w.Nodes = append(w.Nodes, from, to)
	w.CreateIntersection(from.ID)
	w.CreateIntersection(to.ID)

	cmd := &CreateRoadCommand{From: from, To: to, MaxSpeed: 40}
	if err := ex.Execute(cmd); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	select {
		case v := <-ch:
			if _, ok := v.(events.RoadCreatedEvent); !ok {
			t.Fatalf("expected RoadCreatedEvent, got %T", v)
		}
		case <-time.After(1 * time.Second):
			t.Fatalf("expected event but none emitted")
	}
}
