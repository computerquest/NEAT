package main

//todo finish
type Neat struct {
	species int
	nps int //networks per species
	connectMutate float64
	nodeMutate float64
	innovation int
	network [][]Network //stores networks in species
	connectionInnovation [][]int //stores innovation number and connection to and from ex: 1, fromNode:2, toNode: 5
}
