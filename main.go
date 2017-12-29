package main

import (
	"fmt"
)

func main() {
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

	neat := GetNeatInstance(100, 2, 1, .3)
	neat.initialize()

	winner := neat.start(data, 20, 50)

	//neat.printNeat()

	fmt.Println()

	printNetwork(&winner)
	fmt.Println("best ", winner.fitness, "error", 1/winner.fitness)
	fmt.Println("result: ", winner.Process(data[0][0]), winner.Process(data[1][0]), winner.Process(data[2][0]), winner.Process(data[3][0])) //1 1 0 0
	fmt.Println("finsihed")
}
