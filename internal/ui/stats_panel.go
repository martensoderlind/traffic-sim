package ui

import (
	"fmt"
	"image/color"
	"traffic-sim/internal/events"
	"traffic-sim/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
)

type StatsPanel struct {
	X, Y          float64
	Width, Height float64
	shadowOffset  float64
	
	bgColor     color.RGBA
	textColor   color.RGBA
	shadowColor color.RGBA
	
	labels []*Label
	
	roadCount      int
	nodeCount      int
	vehicleCount   int
	spawnCount     int
	despawnCount   int
	trafficLights  int
	
	world        *world.World
	unsubscribers []func()
}

func NewStatsPanel(x, y float64, w *world.World) *StatsPanel {
	panel := &StatsPanel{
		X:            x,
		Y:            y,
		Width:        280,
		Height:       240,
		shadowOffset: 3,
		bgColor:      color.RGBA{40, 40, 50, 240},
		textColor:    color.RGBA{220, 220, 220, 255},
		shadowColor:  color.RGBA{0, 0, 0, 80},
		labels:       make([]*Label, 0),
		world:        w,
		unsubscribers: make([]func(), 0),
	}
	
	panel.initializeStats()
	panel.setupUI()
	panel.subscribeToEvents()
	
	return panel
}

func (p *StatsPanel) initializeStats() {
	p.world.Mu.RLock()
	defer p.world.Mu.RUnlock()
	
	p.roadCount = len(p.world.Roads)
	p.nodeCount = len(p.world.Nodes)
	p.vehicleCount = len(p.world.Vehicles)
	p.spawnCount = len(p.world.SpawnPoints)
	p.despawnCount = len(p.world.DespawnPoints)
	p.trafficLights = len(p.world.TrafficLights)
}

func (p *StatsPanel) setupUI() {
	yOffset := p.Y + 15
	
	titleLabel := NewLabel(p.X+15, yOffset, "Simulation Stats")
	titleLabel.Size = 16
	titleLabel.Color = color.RGBA{255, 255, 255, 255}
	p.labels = append(p.labels, titleLabel)
	yOffset += 30
	
	p.labels = append(p.labels, NewLabel(p.X+15, yOffset, fmt.Sprintf("Roads: %d", p.roadCount)))
	yOffset += 25
	
	p.labels = append(p.labels, NewLabel(p.X+15, yOffset, fmt.Sprintf("Nodes: %d", p.nodeCount)))
	yOffset += 25
	
	p.labels = append(p.labels, NewLabel(p.X+15, yOffset, fmt.Sprintf("Vehicles: %d", p.vehicleCount)))
	yOffset += 25
	
	p.labels = append(p.labels, NewLabel(p.X+15, yOffset, fmt.Sprintf("Spawn Points: %d", p.spawnCount)))
	yOffset += 25
	
	p.labels = append(p.labels, NewLabel(p.X+15, yOffset, fmt.Sprintf("Despawn Points: %d", p.despawnCount)))
	yOffset += 25
	
	p.labels = append(p.labels, NewLabel(p.X+15, yOffset, fmt.Sprintf("Traffic Lights: %d", p.trafficLights)))
	
	for _, label := range p.labels {
		label.Size = 14
	}
}

func (p *StatsPanel) subscribeToEvents() {
	if p.world.Events == nil {
		return
	}
	
	unsub1 := p.world.Events.Subscribe(events.EventRoadCreated, func(payload any) {
		p.roadCount++
		p.updateLabels()
	})
	p.unsubscribers = append(p.unsubscribers, unsub1)
	
	unsub2 := p.world.Events.Subscribe(events.EventRoadDeleted, func(payload any) {
		p.roadCount--
		p.updateLabels()
	})
	p.unsubscribers = append(p.unsubscribers, unsub2)
	
	unsub3 := p.world.Events.Subscribe(events.EventNodeCreated, func(payload any) {
		p.nodeCount++
		p.updateLabels()
	})
	p.unsubscribers = append(p.unsubscribers, unsub3)
	
	unsub4 := p.world.Events.Subscribe(events.EventNodeDeleted, func(payload any) {
		p.nodeCount--
		p.updateLabels()
	})
	p.unsubscribers = append(p.unsubscribers, unsub4)
	
	unsub5 := p.world.Events.Subscribe(events.EventSpawnPointCreated, func(payload any) {
		p.spawnCount++
		p.updateLabels()
	})
	p.unsubscribers = append(p.unsubscribers, unsub5)
	
	unsub6 := p.world.Events.Subscribe(events.EventDespawnPointCreated, func(payload any) {
		p.despawnCount++
		p.updateLabels()
	})
	p.unsubscribers = append(p.unsubscribers, unsub6)
	
	unsub7 := p.world.Events.Subscribe(events.EventTrafficLightCreated, func(payload any) {
		p.trafficLights++
		p.updateLabels()
	})
	p.unsubscribers = append(p.unsubscribers, unsub7)
	
	unsub8 := p.world.Events.Subscribe(events.EventWorldLoaded, func(payload any) {
		ev, ok := payload.(events.WorldLoadedEvent)
		if !ok {
			return
		}
		newWorld := ev.World.(*world.World)
		p.ReplaceWorld(newWorld)
	})
	p.unsubscribers = append(p.unsubscribers, unsub8)
}

func (p *StatsPanel) updateLabels() {
	if len(p.labels) < 7 {
		return
	}
	
	p.labels[1].Text = fmt.Sprintf("Roads: %d", p.roadCount)
	p.labels[2].Text = fmt.Sprintf("Nodes: %d", p.nodeCount)
	p.labels[3].Text = fmt.Sprintf("Vehicles: %d", p.vehicleCount)
	p.labels[4].Text = fmt.Sprintf("Spawn Points: %d", p.spawnCount)
	p.labels[5].Text = fmt.Sprintf("Despawn Points: %d", p.despawnCount)
	p.labels[6].Text = fmt.Sprintf("Traffic Lights: %d", p.trafficLights)
}

func (p *StatsPanel) Update() {
	p.world.Mu.RLock()
	currentVehicleCount := len(p.world.Vehicles)
	p.world.Mu.RUnlock()
	
	if currentVehicleCount != p.vehicleCount {
		p.vehicleCount = currentVehicleCount
		p.updateLabels()
	}
}

func (p *StatsPanel) Draw(screen *ebiten.Image) {
	NewRect(
		float32(p.X+p.shadowOffset),
		float32(p.Y+p.shadowOffset),
		float32(p.Width),
		float32(p.Height),
		13,
		p.shadowColor,
	).draw(screen)
	
	NewRect(
		float32(p.X),
		float32(p.Y),
		float32(p.Width),
		float32(p.Height),
		10,
		p.bgColor,
	).draw(screen)
	
	for _, label := range p.labels {
		label.Draw(screen)
	}
}

func (p *StatsPanel) ReplaceWorld(newWorld *world.World) {
	for _, unsub := range p.unsubscribers {
		unsub()
	}
	p.unsubscribers = make([]func(), 0)
	
	p.world = newWorld
	p.initializeStats()
	p.updateLabels()
	p.subscribeToEvents()
}

func (p *StatsPanel) Cleanup() {
	for _, unsub := range p.unsubscribers {
		unsub()
	}
	p.unsubscribers = nil
}