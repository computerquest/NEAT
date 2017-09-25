package main

import ("fmt")

//todo make sure no errors in expanding the length of the slice (will have issues)
//todo change the append stuff
func main() {
	//n := GetNetworkInstance(10, 10, 0)
	//inputValues  := []float64{.1,.2,.3,.4,.5,.6,.7,.8,.9,.10}
	//n.Process(inputValues)

	a := make([]int, 5, 8)
	for i := 0; i < 10; i++ {
		a = append(a, i)
		fmt.Println("%o\n", len(a), cap(a), a[0], a)
	}

	fmt.Print("works")
}
