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
		Position: 0,
		Speed: 20,
		MaxSpeed:25,
	}

	world := sim.NewWorld([]*road.Road{r1,r2,r3,r4},[]*road.Node{n1,n2,n3},[]*vehicle.Vehicle{v1})
	simulator := sim.NewSimulator(world, 100*time.Millisecond)

	// intersections := sim.BuildIntersections([]*road.Road{r1,r2,r3,r4}, []*road.Node{n1, n2,n3})
	// fmt.Println("Intersections:", len(intersections))
	// for  _,ints:=range intersections{
	// 	fmt.Println("intersection connections:",ints.ID)
	// 	fmt.Println("intersection incoming:",ints.Incoming)
	// 	fmt.Println("intersection outgoing:",ints.Outgoing)
	// }

	fmt.Println("starting simulation..")
	go simulator.Start()

	for range time.Tick(time.Second){
		world.Mu.RLock()
		fmt.Printf("%s: pos=%.2f on %s\n", v1.ID, v1.Position, v1.Road.ID)
		world.Mu.RUnlock()
	}
}