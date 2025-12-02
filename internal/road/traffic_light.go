package road

type LightState int

const (
	LightRed LightState = iota
	LightYellow
	LightGreen
)

type TrafficLight struct {
	ID           string
	Intersection *Intersection
	ControlledRoads []*Road
	PrevState	LightState
	State        LightState
	Timer        float64
	GreenTime    float64
	YellowTime   float64
	RedTime      float64
	Enabled      bool
}
func NewTrafficLight(id string, intersection *Intersection, startGreen bool) *TrafficLight {
	initialState := LightRed
	if startGreen {
		initialState = LightGreen
	}
	
	return &TrafficLight{
		ID:           id,
		Intersection: intersection,
		ControlledRoads: make([]*Road, 0),
		State:        initialState,
		Timer:        0.0,
		GreenTime:    8.0,
		YellowTime:   2.0,
		RedTime:      8.0,
		Enabled:      true,
	}
}

func (tl *TrafficLight) AddControlledRoad(r *Road) {
	tl.ControlledRoads = append(tl.ControlledRoads, r)
}

func (tl *TrafficLight) Update(dt float64) {
	if !tl.Enabled {
		return
	}

	tl.Timer += dt

	switch tl.State {
	case LightGreen:
		if tl.Timer >= tl.GreenTime {
			tl.PrevState = tl.State
			tl.State = LightYellow
			tl.Timer = 0.0
		}
	case LightYellow:
		if tl.Timer >= tl.YellowTime {
			if tl.PrevState == LightGreen {
				tl.State = LightRed
			} else {
				tl.State = LightGreen
			}
			tl.Timer = 0.0
		}
	case LightRed:
		if tl.Timer >= tl.RedTime {
			tl.PrevState = tl.State
			tl.State = LightYellow
			tl.Timer = 0.0
		}
	}
}

func (tl *TrafficLight) CanProceed() bool {
	return tl.State == LightGreen || !tl.Enabled
}

func (tl *TrafficLight) ShouldSlow() bool {
	return tl.State == LightYellow && tl.Enabled
}

func (tl *TrafficLight) IsRed() bool {
	return tl.State == LightRed && tl.Enabled
}