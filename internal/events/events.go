package events

const (
    EventRoadCreated = "road.created"
    EventWorldLoaded  = "world.loaded"
)

type RoadCreatedEvent struct {
    Road any
}

type WorldLoadedEvent struct {
    World any
}
