package main

import ("fmt")

func main() {
	//xor
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
			{1},
		},
	}

	network := GetNetworkInstance(2, 1, 1, 1)
	network.trainSet(data)


	fmt.Print("works")
}
