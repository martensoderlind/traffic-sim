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
	n2:= &road.Node{ID:"n1",X:100,Y:0}

	r1 :=road.NewRoad("r1",n1,n2,20)

	v1 := &vehicle.Vehicle{
		ID:"car1",
		Road: r1,
		Position: 0,
		Speed: 10,
		MaxSpeed:15,
	}

	world := sim.NewWorld([]*road.Road{r1},[]*road.Node{n1,n2},[]*vehicle.Vehicle{v1})
	simulator := sim.NewSimulation(world, 100*time.Millisecond)

	fmt.Println("starting simulation..")
	go simulator.Start()

	for range time.Tick(time.Second){
		world.Mu.RLock()
		fmt.Printf("%s: pos=%.2f on %s\n", v1.ID, v1.Position, v1.Road.ID)
		world.Mu.RUnlock()
	}
}