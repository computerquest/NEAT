package main

import (
	"fmt"
)

func main() {
	//XOR data seti
	data := [][][]float64{
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

	var winner Network
	for i := 0; i < 1000; i++ {
		neat := GetNeatInstance(10, 2, 1, .3, .1)
		neat.initialize()

		winner = neat.start(data, 20, 50)
	}
	//neat.printNeat()

	fmt.Println()

	printNetwork(&winner)
	fmt.Println("best ", winner.fitness, "error", 1/winner.fitness)
	fmt.Println("result: ", winner.Process(data[0][0]), winner.Process(data[1][0]), winner.Process(data[2][0]), winner.Process(data[3][0])) //1 1 0 0
	fmt.Println("finsihed")
}
