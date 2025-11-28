package sim

import (
	"math/rand"
	"sync"
	"time"

	"traffic-sim/internal/road"
	"traffic-sim/internal/vehicle"
)

type Simulator struct {
	world *World
	tickRate time.Duration
}
type World struct {
 	Roads         []*road.Road
    Nodes         []*road.Node
    Vehicles      []*vehicle.Vehicle
    Intersections []*road.Intersection

    // Fast lookup
    IntersectionsByNode map[string]*road.Intersection

    Mu sync.RWMutex
}

func BuildIntersections(roads []*road.Road, nodes []*road.Node) []*road.Intersection {
	m := make(map[string]*road.Intersection)

	for _, n := range nodes {
		m[n.ID] = road.NewIntersection(n.ID)
	}

	for _, r := range roads {
		in := m[r.From.ID]
		out := m[r.To.ID]
		in.AddOutgoing(r)
		out.AddIncoming(r)
	}

	intersections := make([]*road.Intersection, 0, len(m))
	for _, i := range m {
		intersections = append(intersections, i)
	}
	return intersections
}

func NewWorld(roads []*road.Road, nodes []*road.Node, vehicles []*vehicle.Vehicle) *World {
    w := &World{
        Roads:               roads,
        Nodes:               nodes,
        Vehicles:            vehicles,
        IntersectionsByNode: make(map[string]*road.Intersection),
    }

    w.Intersections = BuildIntersections(roads, nodes)
    for _, i := range w.Intersections {
        w.IntersectionsByNode[i.ID] = i
    }

    return w
}

func NewSimulator(world *World, tickRate time.Duration) *Simulator {
	return &Simulator{world: world, tickRate: tickRate}
}

func (s *Simulator) Start() {
	ticker := time.NewTicker(s.tickRate)
	for range ticker.C {
		s.update()
	}
}

func (s *Simulator) nextRoadFor(v *vehicle.Vehicle) *road.Road {
    for _, i := range s.world.Intersections {
        if i.ID == v.Road.To.ID {
            if len(i.Outgoing) == 0 {
                return nil
            }
            
            available := make([]*road.Road, 0, len(i.Outgoing))
            for _, r := range i.Outgoing {
                if r.ID != v.Road.ID { 
                    available = append(available, r)
                }
            }
            
            if len(available) == 0 {
                return nil  
            }
            
            return available[rand.Intn(len(available))]
        }
    }
    return nil
}

func (s *Simulator) update() {
	s.world.Mu.Lock()
	defer s.world.Mu.Unlock()

	dt := s.tickRate.Seconds()

	for _, v := range s.world.Vehicles {
		newPos := v.Position + v.Speed*dt
		if newPos >= v.Road.Length {
			next := s.nextRoadFor(v)
		if next != nil {
			v.Road = next
			v.Position = 0
		} else {
			v.Position = v.Road.Length
			v.Speed = 0
		}
		} else {
			v.Position = newPos
		}
	}
}