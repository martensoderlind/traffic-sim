package persistence

const CurrentVersion = "1.0.0"

type SaveFormat struct {
	Version       string                 `json:"version"`
	Timestamp     string                 `json:"timestamp"`
	Nodes         []NodeData             `json:"nodes"`
	Roads         []RoadData             `json:"roads"`
	SpawnPoints   []SpawnPointData       `json:"spawnPoints"`
	DespawnPoints []DespawnPointData     `json:"despawnPoints"`
	TrafficLights []TrafficLightData     `json:"trafficLights"`
}

type NodeData struct {
	ID string  `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}

type RoadData struct {
	ID              string  `json:"id"`
	FromNodeID      string  `json:"fromNodeId"`
	ToNodeID        string  `json:"toNodeId"`
	MaxSpeed        float64 `json:"maxSpeed"`
	Width           float64 `json:"width"`
	ReverseRoadID   string  `json:"reverseRoadId,omitempty"`
	StartOffsetX    float64 `json:"startOffsetX,omitempty"`
	StartOffsetY    float64 `json:"startOffsetY,omitempty"`
	EndOffsetX      float64 `json:"endOffsetX,omitempty"`
	EndOffsetY      float64 `json:"endOffsetY,omitempty"`
}

type SpawnPointData struct {
	ID             string  `json:"id"`
	NodeID         string  `json:"nodeId"`
	RoadID         string  `json:"roadId"`
	Interval       float64 `json:"interval"`
	MinSpeed       float64 `json:"minSpeed"`
	MaxSpeed       float64 `json:"maxSpeed"`
	MaxVehicles    int     `json:"maxVehicles"`
	Enabled        bool    `json:"enabled"`
	VehicleCounter int     `json:"vehicleCounter"`
}

type DespawnPointData struct {
	ID      string `json:"id"`
	NodeID  string `json:"nodeId"`
	RoadID  string `json:"roadId"`
	Enabled bool   `json:"enabled"`
}

type TrafficLightData struct {
	ID                string   `json:"id"`
	IntersectionID    string   `json:"intersectionId"`
	ControlledRoadIDs []string `json:"controlledRoadIds"`
	State             int      `json:"state"`
	GreenTime         float64  `json:"greenTime"`
	YellowTime        float64  `json:"yellowTime"`
	RedTime           float64  `json:"redTime"`
	Enabled           bool     `json:"enabled"`
}