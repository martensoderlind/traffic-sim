package sim

import (
	"sync"
	"time"

	"traffic-sim/internal/road"
	"traffic-sim/internal/vehicle"
)

type World struct {
	Roads []*road.Road
	Nodes []*road.Node
	Vehicles []*vehicle.Vehicle

	Mu sync.RWMutex
}

func NewWorld(roads []*road.Road, nodes []*road.Node, vehicles []*vehicle.Vehicle) *World{

	return &World{
		Roads: roads,
		Nodes:nodes,
		Vehicles:vehicles,
	}
}


type Simulator struct{
	World *World
	TickRate time.Duration
}

func NewSimulation(world *World, tickRate time.Duration) *Simulator{
	return &Simulator{World: world,TickRate: tickRate}
}

func (s *Simulator) Start(){
	ticker:= time.NewTicker(s.TickRate)
	for range ticker.C{
		s.update()
	}
}

func (s *Simulator) update(){
	s.World.Mu.Lock()
	defer s.World.Mu.Unlock()

	dt:= s.TickRate.Seconds()

	for _,v :=range s.World.Vehicles {
		newPos := v.Position +v.Speed*dt
		if newPos >=v.Road.Length{
			v.Position = v.Road.Length
			v.Speed = 0
		}else{
			v.Position = newPos
		}
	}
}

func BuildIntersections(roads []*road.Road, nodes []*road.Node) []*road.Intersection {
m := make(map[string]*road.Intersection)


for _, n := range nodes {
m[n.ID] = road.NewIntersection(n.ID)
}


// Assign roads to intersections
for _, r := range roads {
in := m[r.From.ID]
out := m[r.To.ID]
in.AddOutgoing(r)
out.AddIncoming(r)
}


// Flatten to slice
intersections := make([]*road.Intersection, 0, len(m))
for _, i := range m {
intersections = append(intersections, i)
}
return intersections
}