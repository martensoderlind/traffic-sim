package renderer

import (
	"image/color"
	"math"
	"traffic-sim/internal/vehicle"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type VehicleRenderer struct{
	shadowOffset float32
	showTargetLines bool
}

func NewVehicleRenderer() *VehicleRenderer {
	return &VehicleRenderer{
		shadowOffset: 2.0,
		showTargetLines: true,
	}
}

func (vr *VehicleRenderer) SetShowTargetLines(show bool) {
	vr.showTargetLines = show
}

func (vr *VehicleRenderer) RenderVehicles(screen *ebiten.Image, vehicles []*vehicle.Vehicle) {
	if vr.showTargetLines {
		for _, v := range vehicles {
			vr.renderTargetLine(screen, v)
		}
	}
	
	for _, v := range vehicles {
		vr.renderSingleVehicle(screen, v)
	}
}

func (vr *VehicleRenderer) renderTargetLine(screen *ebiten.Image, v *vehicle.Vehicle) {
	if v.TargetDespawn == nil {
		return
	}
	
	targetX := float32(v.TargetDespawn.Node.X)
	targetY := float32(v.TargetDespawn.Node.Y)
	vehicleX := float32(v.Pos.X)
	vehicleY := float32(v.Pos.Y)
	
	lineColor := color.RGBA{100, 100, 255, 80}
	vector.StrokeLine(screen, vehicleX, vehicleY, targetX, targetY, 1, lineColor, false)
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
	shadow := make([][2]float32, 4)
	
	for i, corner := range corners {
		rotated[i][0] = cx + corner[0]*cos - corner[1]*sin
		rotated[i][1] = cy + corner[0]*sin + corner[1]*cos
		
		shadow[i][0] = rotated[i][0] + vr.shadowOffset
		shadow[i][1] = rotated[i][1] + vr.shadowOffset
	}

	shadowVerts := []ebiten.Vertex{
		{DstX: shadow[0][0], DstY: shadow[0][1], SrcX: 0, SrcY: 0, ColorR: 0, ColorG: 0, ColorB: 0, ColorA: 0.4},
		{DstX: shadow[1][0], DstY: shadow[1][1], SrcX: 0, SrcY: 0, ColorR: 0, ColorG: 0, ColorB: 0, ColorA: 0.4},
		{DstX: shadow[2][0], DstY: shadow[2][1], SrcX: 0, SrcY: 0, ColorR: 0, ColorG: 0, ColorB: 0, ColorA: 0.4},
		{DstX: shadow[3][0], DstY: shadow[3][1], SrcX: 0, SrcY: 0, ColorR: 0, ColorG: 0, ColorB: 0, ColorA: 0.4},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(shadowVerts, indices, vr.createWhiteImage(), nil)

	var bodyColor1, bodyColor2 color.RGBA
	if v.TargetDespawn != nil {
		bodyColor1 = color.RGBA{90, 90, 230, 255}
		bodyColor2 = color.RGBA{70, 70, 180, 255}
	} else {
		bodyColor1 = color.RGBA{230, 65, 65, 255}
		bodyColor2 = color.RGBA{180, 40, 40, 255}
	}

	bodyVerts := []ebiten.Vertex{
		{DstX: rotated[0][0], DstY: rotated[0][1], SrcX: 0, SrcY: 0, 
			ColorR: float32(bodyColor1.R)/255, ColorG: float32(bodyColor1.G)/255, 
			ColorB: float32(bodyColor1.B)/255, ColorA: 1},
		{DstX: rotated[1][0], DstY: rotated[1][1], SrcX: 0, SrcY: 0, 
			ColorR: float32(bodyColor1.R)/255, ColorG: float32(bodyColor1.G)/255, 
			ColorB: float32(bodyColor1.B)/255, ColorA: 1},
		{DstX: rotated[2][0], DstY: rotated[2][1], SrcX: 0, SrcY: 0, 
			ColorR: float32(bodyColor2.R)/255, ColorG: float32(bodyColor2.G)/255, 
			ColorB: float32(bodyColor2.B)/255, ColorA: 1},
		{DstX: rotated[3][0], DstY: rotated[3][1], SrcX: 0, SrcY: 0, 
			ColorR: float32(bodyColor2.R)/255, ColorG: float32(bodyColor2.G)/255, 
			ColorB: float32(bodyColor2.B)/255, ColorA: 1},
	}
	screen.DrawTriangles(bodyVerts, indices, vr.createWhiteImage(), nil)

	windshieldH := height * 0.25
	windshield := [][2]float32{
		{-hw * 0.7, -hh},
		{hw * 0.7, -hh},
		{hw * 0.5, -hh + windshieldH},
		{-hw * 0.5, -hh + windshieldH},
	}
	
	windRotated := make([][2]float32, 4)
	for i, corner := range windshield {
		windRotated[i][0] = cx + corner[0]*cos - corner[1]*sin
		windRotated[i][1] = cy + corner[0]*sin + corner[1]*cos
	}

	windVerts := []ebiten.Vertex{
		{DstX: windRotated[0][0], DstY: windRotated[0][1], SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.5, ColorB: 0.6, ColorA: 0.8},
		{DstX: windRotated[1][0], DstY: windRotated[1][1], SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.5, ColorB: 0.6, ColorA: 0.8},
		{DstX: windRotated[2][0], DstY: windRotated[2][1], SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.4, ColorB: 0.5, ColorA: 0.6},
		{DstX: windRotated[3][0], DstY: windRotated[3][1], SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.4, ColorB: 0.5, ColorA: 0.6},
	}
	screen.DrawTriangles(windVerts, indices, vr.createWhiteImage(), nil)

	headlightOffset := hh * 0.8
	headlightSize := float32(1.5)
	
	leftHeadlight := [2]float32{
		cx + (-hw*0.5)*cos - (-headlightOffset)*sin,
		cy + (-hw*0.5)*sin + (-headlightOffset)*cos,
	}
	rightHeadlight := [2]float32{
		cx + (hw*0.5)*cos - (-headlightOffset)*sin,
		cy + (hw*0.5)*sin + (-headlightOffset)*cos,
	}

	vector.FillCircle(screen, leftHeadlight[0], leftHeadlight[1], headlightSize, 
		color.RGBA{255, 255, 200, 255}, false)
	vector.FillCircle(screen, rightHeadlight[0], rightHeadlight[1], headlightSize, 
		color.RGBA{255, 255, 200, 255}, false)

	if v.Speed < 1.0 {
		tailOffset := hh * 0.9
		tailSize := float32(1.2)
		
		leftTail := [2]float32{
			cx + (-hw*0.6)*cos - (tailOffset)*sin,
			cy + (-hw*0.6)*sin + (tailOffset)*cos,
		}
		rightTail := [2]float32{
			cx + (hw*0.6)*cos - (tailOffset)*sin,
			cy + (hw*0.6)*sin + (tailOffset)*cos,
		}

		vector.FillCircle(screen, leftTail[0], leftTail[1], tailSize, 
			color.RGBA{255, 50, 50, 255}, false)
		vector.FillCircle(screen, rightTail[0], rightTail[1], tailSize, 
			color.RGBA{255, 50, 50, 255}, false)
	}

	var edgeColor color.RGBA
	if v.TargetDespawn != nil {
		edgeColor = color.RGBA{70, 70, 180, 255}
	} else {
		edgeColor = color.RGBA{180, 50, 50, 255}
	}
	
	vector.StrokeLine(screen, rotated[0][0], rotated[0][1], rotated[1][0], rotated[1][1], 1, edgeColor, false)
	vector.StrokeLine(screen, rotated[1][0], rotated[1][1], rotated[2][0], rotated[2][1], 1, edgeColor, false)
	vector.StrokeLine(screen, rotated[2][0], rotated[2][1], rotated[3][0], rotated[3][1], 1, edgeColor, false)
	vector.StrokeLine(screen, rotated[3][0], rotated[3][1], rotated[0][0], rotated[0][1], 1, edgeColor, false)
}

func (vr *VehicleRenderer) createWhiteImage() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}