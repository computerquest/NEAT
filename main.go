package main

import (
	"fmt"
)

func main() {
	//xor
	/*data := [][][]float64{
		{
			{0, 1},
			{1},
		},
		{
			{1, 0},
			{1},
		},
		{
			{0, 0},
			{0},
		},
		{
			{1, 1},
			{0},
		},
	}

	dataA := [][][]float64{
		{
			{.05, .1},
			{.01, .99},
		},
	}

	cA := make(chan float64)
	//dataA configuration
	networkA := GetNetworkInstance(2, 2, 1, 1, .01)
	networkA.mutateNode(2, 0, 10, 11)
	networkA.mutateNode(2, 1, 10, 11)
	networkA.mutateNode(3, 0, 12, 13)
	networkA.mutateNode(3, 1, 12, 13)
	stuff := func() {
		networkA.trainSet(dataA, 10000)
		cA <- 1
	}
	go stuff()

	fmt.Println()

	c := make(chan float64)

	//data configuration
	network := GetNetworkInstance(2, 1, 1, 1, .01)
	network.mutateNode(2, 0, 10, 11)
	network.mutateNode(1, 0, 12, 13)
	network.mutateConnection(1, 4, 15)
	network.mutateConnection(2, 3, 16)
	stuffA := func() {
		network.trainSet(data, 10000)
		c <- 1
	}
	go stuffA()

	<-cA
	<-c

	fmt.Println(networkA.Process(dataA[0][0]))
	fmt.Println(network.Process(data[0][0]), network.Process(data[1][0]), network.Process(data[2][0]))

	fmt.Print("works")*/

	dataA := [][][]float64{
		{
			{.05, .1},
			{.01, .99},
		},
	}

	neat := GetNeatInstance(15, 2, 2)

	for i := 1000; i >= 0; i-- {
		neat.start(dataA)

		if i%5 == 0 {
			neat.speciateAll()
		}
	}

	neat.printNeat()
	fmt.Println("finsihed")
}
