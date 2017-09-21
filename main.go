package main

import ("fmt")

//todo make sure no errors in expanding the length of the slice (will have issues)
func main() {
	n := GetNetworkInstance(10, 10, 0)
	inputValues  := []float64{.1,.2,.3,.4,.5,.6,.7,.8,.9,.10}
	n.Process(inputValues)
	fmt.Print("works")
}
