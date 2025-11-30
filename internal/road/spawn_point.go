package road

type SpawnPoint struct {
	ID            string
	Node          *Node
	Road          *Road
	Interval      float64
	Timer         float64
	MinSpeed      float64
	MaxSpeed      float64
	MaxVehicles   int
	Enabled       bool
	VehicleCounter int
}

func NewSpawnPoint(id string, node *Node, road *Road) *SpawnPoint {
	return &SpawnPoint{
		ID:            id,
		Node:          node,
		Road:          road,
		Interval:      3.0,
		Timer:         0.0,
		MinSpeed:      20.0,
		MaxSpeed:      40.0,
		MaxVehicles:   50,
		Enabled:       true,
		VehicleCounter: 0,
	}
}