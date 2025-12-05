package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"traffic-sim/internal/input"
	"traffic-sim/internal/renderer"
	"traffic-sim/internal/sim"
	"traffic-sim/internal/world"
)

type Game struct {
	renderer     *renderer.Renderer
	simulator    *sim.Simulator
	world        *world.World
	InputHandler *input.InputHandler
	lastTime     time.Time
}

func (g *Game) Update() error {
	mouseX, mouseY := ebiten.CursorPosition()
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	
	g.renderer.Toolbar.Update(mouseX, mouseY, clicked)
	g.InputHandler.Update()
	
	now := time.Now()
	if !g.lastTime.IsZero() {
		dt := now.Sub(g.lastTime).Seconds()
		g.simulator.UpdateOnce(dt)
	}
	g.lastTime = now
	
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.renderer.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.renderer.Layout(outsideWidth, outsideHeight)
}

func (g *Game) replaceWorld(newWorld *world.World) {
	g.world.Mu.Lock()
	g.world = newWorld
	g.world.Mu.Unlock()

	g.simulator = sim.NewSimulator(g.world, 8*time.Millisecond)
	
	g.InputHandler.ReplaceWorld(g.world)
	g.InputHandler.Simulator = g.simulator
	
	g.renderer.World = g.world
	g.renderer.InputHandler = g.InputHandler
	
	
	log.Println("World replaced successfully")
}

func main() {
	world := world.New()
	simulator := sim.NewSimulator(world, 8*time.Millisecond)
	inputHandler := input.NewInputHandler(world, simulator)

	rend := renderer.NewRenderer(world, inputHandler)

	game := &Game{
		renderer:     rend,
		simulator:    simulator,
		world:        world,
		InputHandler: inputHandler,
	}
	
	inputHandler.SetOnWorldReplaced(game.replaceWorld)

	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("Traffic Simulation")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
