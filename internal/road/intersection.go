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