package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"traffic-sim/internal/input"
	"traffic-sim/internal/renderer"
	"traffic-sim/internal/road"
	"traffic-sim/internal/sim"
	"traffic-sim/internal/vehicle"
)

type Game struct {
	renderer  *renderer.Renderer
	simulator *sim.Simulator
	world     *sim.World
	InputHandler *input.InputHandler
}

func (g *Game) Update() error {
	g.InputHandler.Update()
	g.simulator.UpdateOnce()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.renderer.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.renderer.Layout(outsideWidth, outsideHeight)
}

func main() {
	n1 := &road.Node{ID: "n1", X: 100, Y: 300}
	n2 := &road.Node{ID: "n2", X: 400, Y: 300}
	n3 := &road.Node{ID: "n3", X: 400, Y: 100}
	n4 := &road.Node{ID: "n4", X: 700, Y: 100}

	r1 := road.NewRoad("r1", n1, n2, 40)
	r2 := road.NewRoad("r2", n2, n3, 40)
	r3 := road.NewRoad("r3", n3, n4, 40)
	r4 := road.NewRoad("r4", n3, n2, 40)
	r5 := road.NewRoad("r5", n2, n1, 40)

	v1 := &vehicle.Vehicle{
		ID:       "car1",
		Road:     r1,
		Distance: 0,
		Pos:      vehicle.Vec2{X: r1.From.X, Y: r1.From.Y},
		Speed:    30,
	}

	v2 := &vehicle.Vehicle{
		ID:       "car2",
		Road:     r2,
		Distance: 0,
		Pos:      vehicle.Vec2{X: r2.From.X, Y: r2.From.Y},
		Speed:    25,
	}

	v3 := &vehicle.Vehicle{
		ID:       "car3",
		Road:     r5,
		Distance: 50,
		Pos:      vehicle.Vec2{X: r5.From.X, Y: r5.From.Y},
		Speed:    35,
	}

	world := sim.NewWorld(
		[]*road.Road{r1, r2, r3, r4, r5},
		[]*road.Node{n1, n2, n3, n4},
		[]*vehicle.Vehicle{v1, v2, v3},
	)

	simulator := sim.NewSimulator(world, 16*time.Millisecond) // ~60 FPS

	inputHandler := input.NewInputHandler(world)

	rend := &renderer.Renderer{World: world}

	game := &Game{
		renderer:  rend,
		simulator: simulator,
		world:     world,
		InputHandler: inputHandler,
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Traffic Simulation")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}