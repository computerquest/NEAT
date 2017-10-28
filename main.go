package main

import ("fmt")

func main() {
	n := GetNetworkInstance(5,5, 0)

	n.mutateNode(5, 0,1,9)

	fmt.Print("works")
}
