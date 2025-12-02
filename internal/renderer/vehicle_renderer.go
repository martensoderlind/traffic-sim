package renderer

import (
	"image/color"
	"math"
	"traffic-sim/internal/vehicle"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type VehicleRenderer struct{}

func NewVehicleRenderer() *VehicleRenderer {
	return &VehicleRenderer{}
}

func (vr *VehicleRenderer) RenderVehicles(screen *ebiten.Image, vehicles []*vehicle.Vehicle) {
	for _, v := range vehicles {
		vr.renderSingleVehicle(screen, v)
	}
}

func (vr *VehicleRenderer) renderSingleVehicle(screen *ebiten.Image, v *vehicle.Vehicle) {
	pos := v.Position()
	angle := v.GetAngle()

	width := float32(5.0)
	height := float32(10.0)

	cx := float32(pos.X)
	cy := float32(pos.Y)

	cos := float32(math.Cos(angle))
	sin := float32(math.Sin(angle))

	hw := width / 2
	hh := height / 2

	corners := [][2]float32{
		{-hw, -hh},
		{hw, -hh},
		{hw, hh},
		{-hw, hh},
	}

	rotated := make([][2]float32, 4)
	for i, corner := range corners {
		rotated[i][0] = cx + corner[0]*cos - corner[1]*sin
		rotated[i][1] = cy + corner[0]*sin + corner[1]*cos
	}

	vertices := []ebiten.Vertex{
		{DstX: rotated[0][0], DstY: rotated[0][1], SrcX: 0, SrcY: 0, ColorR: 1.0, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		{DstX: rotated[1][0], DstY: rotated[1][1], SrcX: 0, SrcY: 0, ColorR: 1.0, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		{DstX: rotated[2][0], DstY: rotated[2][1], SrcX: 0, SrcY: 0, ColorR: 1.0, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		{DstX: rotated[3][0], DstY: rotated[3][1], SrcX: 0, SrcY: 0, ColorR: 1.0, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
	}

	indices := []uint16{0, 1, 2, 0, 2, 3}

	screen.DrawTriangles(vertices, indices, emptyImage(screen), nil)

	vector.StrokeLine(screen, rotated[0][0], rotated[0][1], rotated[1][0], rotated[1][1], 1, color.RGBA{255, 150, 150, 255}, false)
	vector.StrokeLine(screen, rotated[1][0], rotated[1][1], rotated[2][0], rotated[2][1], 1, color.RGBA{255, 150, 150, 255}, false)
	vector.StrokeLine(screen, rotated[2][0], rotated[2][1], rotated[3][0], rotated[3][1], 1, color.RGBA{255, 150, 150, 255}, false)
	vector.StrokeLine(screen, rotated[3][0], rotated[3][1], rotated[0][0], rotated[0][1], 1, color.RGBA{255, 150, 150, 255}, false)
}