package main

import (
	"fmt"
	"time"
	"traffic-sim/internal/road"
	"traffic-sim/internal/sim"
	"traffic-sim/internal/vehicle"
)

func main(){


	n1:= &road.Node{ID:"n1",X:0,Y:0}
	n2:= &road.Node{ID:"n2",X:100,Y:0}
	n3:= &road.Node{ID:"n3",X:100,Y:100}

	r1 :=road.NewRoad("r1",n1,n2,40)
	r2 :=road.NewRoad("r2",n2,n3,40)
	r3 :=road.NewRoad("r3",n3,n2,40)
	r4 :=road.NewRoad("r4",n2,n1,40)

	v1 := &vehicle.Vehicle{
		ID:"car1",
		Road: r1,
		Distance: 0,
		Pos: vehicle.Vec2{X: r1.From.X,Y: r1.From.Y},
		Speed: 20,
	}

	world := sim.NewWorld([]*road.Road{r1,r2,r3,r4},[]*road.Node{n1,n2,n3},[]*vehicle.Vehicle{v1})
	simulator := sim.NewSimulator(world, 100*time.Millisecond)

	fmt.Println("starting simulation..")
	go simulator.Start()

	for range time.Tick(time.Second){
		world.Mu.RLock()
		fmt.Printf("%s: pos=%.2f on %s\n", v1.ID, v1.Pos, v1.Road.ID)
		world.Mu.RUnlock()
	}
}