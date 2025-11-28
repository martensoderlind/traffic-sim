package renderer

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"traffic-sim/internal/sim"
)

type Renderer struct {
	World *sim.World
}

func (r *Renderer) Update() error {
	return nil
}

func (r *Renderer) renderVehicles(screen *ebiten.Image){
		for _, v := range r.World.Vehicles {
		pos := v.Position()
		px := float32(pos.X)
		py := float32(pos.Y)

		vector.FillRect(
			screen,
			px-3, py-3,
			6, 6,
			color.RGBA{255, 50, 50, 255},
			false,
		)

		vector.StrokeRect(
			screen,
			px-3, py-3,
			6, 6,
			1,
			color.RGBA{255, 150, 150, 255},
			false,
		)
	}
}
func (r *Renderer) renderRoads(screen *ebiten.Image){
		for _, rd := range r.World.Roads {
		x1 := float32(rd.From.X)
		y1 := float32(rd.From.Y)
		x2 := float32(rd.To.X)
		y2 := float32(rd.To.Y)

		vector.StrokeLine(
			screen,
			x1, y1,
			x2, y2,
			3,
			color.RGBA{80, 80, 90, 255},
			false,
		)

		vector.StrokeLine(
			screen,
			x1, y1,
			x2, y2,
			1,
			color.RGBA{150, 150, 150, 120},
			false,
		)
	}
}
func (r *Renderer) renderNodes(screen *ebiten.Image){
	for _, node := range r.World.Nodes {
		x := float32(node.X)
		y := float32(node.Y)
		
		vector.FillCircle(
			screen,
			x, y,
			6,
			color.RGBA{100, 100, 120, 255},
			false,
		)
	}
}

func (r *Renderer) Draw(screen *ebiten.Image) {
	r.World.Mu.RLock()
	defer r.World.Mu.RUnlock()

	screen.Fill(color.RGBA{20, 20, 30, 255})
	
	r.renderNodes(screen)
	r.renderRoads(screen)
	r.renderVehicles(screen)
}

func (r *Renderer) Layout(w, h int) (int, int) {
	return 800, 600
}

