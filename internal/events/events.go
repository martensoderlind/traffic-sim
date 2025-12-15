package events

const (
    EventRoadCreated = "road.created"
    EventWorldLoaded  = "world.loaded"
    EventSpawnPointCreated = "spawnpoint.created"
)

type RoadCreatedEvent struct {
    Road any
}
type SpawnPointCreatedEvent struct {
    SpawnPoint any
}

type WorldLoadedEvent struct {
    World any
}
