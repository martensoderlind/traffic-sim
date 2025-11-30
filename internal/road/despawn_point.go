package road

type DespawnPoint struct {
	ID      string
	Node    *Node
	Road    *Road
	Enabled bool
}

func NewDespawnPoint(id string, node *Node, road *Road) *DespawnPoint {
	return &DespawnPoint{
		ID:      id,
		Node:    node,
		Road:    road,
		Enabled: true,
	}
}