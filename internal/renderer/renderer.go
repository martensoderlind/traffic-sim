package renderer

import (
	"image/color"
	"traffic-sim/internal/input"
	"traffic-sim/internal/ui"
	"traffic-sim/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
)

type Renderer struct {
	World        *world.World
	InputHandler *input.InputHandler
	Toolbar      *ui.Toolbar
	screenWidth  int
	screenHeight int

	roadRenderer    *RoadRenderer
	vehicleRenderer *VehicleRenderer
	overlayRenderer *OverlayRenderer
	markerRenderer  *MarkerRenderer
}

func NewRenderer(w *world.World, inputHandler *input.InputHandler) *Renderer {
	return &Renderer{
		World:           w,
		InputHandler:    inputHandler,
		Toolbar:         ui.NewToolbar(inputHandler,w),
		screenWidth:     1920,
		screenHeight:    1080,
		roadRenderer:    NewRoadRenderer(),
		vehicleRenderer: NewVehicleRenderer(),
		overlayRenderer: NewOverlayRenderer(),
		markerRenderer:  NewMarkerRenderer(),
	}
}

func (r *Renderer) Update() error {
	return nil
}

func (r *Renderer) Draw(screen *ebiten.Image) {
	r.World.Mu.RLock()
	defer r.World.Mu.RUnlock()

	screen.Fill(color.RGBA{20, 20, 30, 255})

	r.roadRenderer.RenderRoads(screen, r.World.Roads,r.World.Nodes)
	r.markerRenderer.RenderSpawnPoints(screen, r.World.SpawnPoints)
	r.markerRenderer.RenderDespawnPoints(screen, r.World.DespawnPoints)
	r.vehicleRenderer.RenderVehicles(screen, r.World.Vehicles)
	r.overlayRenderer.RenderToolOverlay(screen, r.InputHandler)
	r.markerRenderer.RenderTrafficLights(screen, r.World.TrafficLights, r.World.Nodes)
	r.Toolbar.Draw(screen)
}

func (r *Renderer) Layout(w, h int) (int, int) {
	r.screenWidth = w
	r.screenHeight = h
	r.Toolbar.UpdatePanelPositions(w, h)
	return w, h
}

func (r *Renderer) ReplaceWorld(newWorld *world.World) {
	r.World = newWorld
	r.Toolbar.ReplaceWorld(newWorld)
}

func (r *Renderer) Cleanup() {
	r.Toolbar.Cleanup()
}