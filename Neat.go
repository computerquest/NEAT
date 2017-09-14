package main

//todo finish
type Neat struct {
	species int //number of species desired
	nps int //networks per species
	connectMutate float64 //odds for connection mutation
	nodeMutate float64 //odds for node mutation
	innovation int //number of innovations
	network [][]Network //stores networks in species
	connectionInnovation [][]int //stores innovation number and connection to and from ex: 1, fromNode:2, toNode: 5
}

func getNeatInstance(numSpeices int, nps int) Neat {
	n := Neat{species: numSpeices, nps: nps, network: make([][]Network, numSpeices)}
	for i := 0; i < len(n.network); i++ {
		n.network[i] = make([]Network, nps)
	}

	return n
}
