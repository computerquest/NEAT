package main

import ("fmt")

//todo make sure no errors in expanding the length of the slice (will have issues)
//todo change the append stuff
func main() {
	n := GetNetworkInstance(5,5, 0)

	n.mutateNode(5, 0,1,9)

	fmt.Print("works")
}
