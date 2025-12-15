package events

import "traffic-sim/internal/road"

const (
	EventRoadCreated          = "road.created"
	EventRoadDeleted          = "road.deleted"
	EventNodeCreated          = "node.created"
	EventNodeDeleted          = "node.deleted"
	EventWorldLoaded          = "world.loaded"
	EventSpawnPointCreated    = "spawnpoint.created"
	EventDespawnPointCreated  = "despawnpoint.created"
	EventTrafficLightCreated  = "trafficlight.created"
	EventRoadPropertiesUpdated = "road.properties.updated"
)

type RoadCreatedEvent struct {
	Road *road.Road
}

type RoadDeletedEvent struct {
	RoadID string
}

type NodeCreatedEvent struct {
	Node *road.Node
}

type NodeDeletedEvent struct {
	NodeID string
}

type SpawnPointCreatedEvent struct {
	SpawnPoint *road.SpawnPoint
}

type DespawnPointCreatedEvent struct {
	DespawnPoint *road.DespawnPoint
}

type TrafficLightCreatedEvent struct {
	TrafficLight *road.TrafficLight
}

type RoadPropertiesUpdatedEvent struct {
	Road     *road.Road
	MaxSpeed float64
	Width    float64
}

type WorldLoadedEvent struct {
	World any
}