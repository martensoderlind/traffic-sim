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
    // Fysik körs i main-loopen av dig själv
    return nil
}

func (r *Renderer) Draw(screen *ebiten.Image) {

    for _, rd := range r.World.Roads {
        x1 := float32(rd.From.X)
        y1 := float32(rd.From.Y)
        x2 := float32(rd.To.X)
        y2 := float32(rd.To.Y)

        vector.StrokeLine(
            screen,
            x1, y1,
            x2, y2,
            2,
            color.White,
            false,
        )
    }

    // Rita fordon – små fyrkanter
    for _, v := range r.World.Vehicles {
        pos := v.Position()
        px := float32(pos.X) - 2
        py := float32(pos.Y) - 2

        vector.FillRect(
            screen,
            px, py,
            4, 4,
            color.RGBA{255, 0, 0, 255},
            false,
        )
    }
}

func (r *Renderer) Layout(w, h int) (int, int) {
    return 800, 600
}
