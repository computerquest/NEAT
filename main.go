package main

import ("fmt")

func main() {
	n := Network{

	}

	n.GetInstance(10, 10)
	inputValues  := []float64{.1,.2,.3,.4,.5,.6,.7,.8,.9,.10}
	n.Process(inputValues)
	fmt.Print("works")
}
