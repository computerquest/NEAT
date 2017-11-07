package main

import ("fmt")

func main() {
	s := GetNeatInstance(40, 5, 5)

	s.mutatePopulation()

	fmt.Print("works")
}
