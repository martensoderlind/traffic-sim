package road

type Intersection struct {
	ID string
	Incoming []*Road
	Outgoing []*Road
}

func NewIntersection(id string) *Intersection {
	return &Intersection{ID: id}
}

func (i *Intersection) AddIncoming(r *Road) {
	i.Incoming = append(i.Incoming, r)
}

func (i *Intersection) AddOutgoing(r *Road) {
	i.Outgoing = append(i.Outgoing, r)
}

func BuildIntersections(roads []*Road, nodes []*Node) []*Intersection {
	m := make(map[string]*Intersection)

	for _, n := range nodes {
		m[n.ID] = NewIntersection(n.ID)
	}

	for _, r := range roads {
		in := m[r.From.ID]
		out := m[r.To.ID]
		in.AddOutgoing(r)
		out.AddIncoming(r)
	}

	intersections := make([]*Intersection, 0, len(m))
	for _, i := range m {
		intersections = append(intersections, i)
	}
	return intersections
}