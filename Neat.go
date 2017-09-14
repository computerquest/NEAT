package main

import (
	"math/rand"
	"time"
	"text/template/parse"
)

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

func GetNeatInstance(numSpeices int, nps int) Neat {
	n := Neat{species: numSpeices, nps: nps, innovation: 0, connectMutate: .7,
		nodeMutate: .3, network: make([][]Network, numSpeices), connectionInnovation: make([][]int, 10)}
	for i := 0; i < len(n.network); i++ {
		n.network[i] = make([]Network, nps)
	}

	for i := 0; i < len(n.connectionInnovation); i++ {
		n.connectionInnovation[i] = make([]int, 2)
	}

	/*REST OF METHOD
	perform initial mutations
		between 1-3 for each network
	speciate
	 */
	return n
}

//todo clarify speciation
/*
are you comparing every network to every other or are you comparing random geneomes (collection of genes) from last generation of species to each network
 */
func (n *Neat) speciate() {
	for i := 0; i < len(n.network); i++ {
		for a := 0; a < len(n.network[i]); a++ {

		}
	}
}
//todo test
func (n *Neat) mutateNetwork(network *Network) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	nodeRange := network.numNodes

	//todo test
	addConnectionInnovation := func(numTo int, numFrom int) int {
		ans := n.innovation
		if len(n.connectionInnovation) >= cap(n.connectionInnovation) {
			newStuff := []int{numFrom, numTo}
			n.connectionInnovation = append(n.connectionInnovation, newStuff)
		} else {
				n.connectionInnovation[n.innovation][0] = numFrom
				n.connectionInnovation[n.innovation][1] = numTo
		}

		network.addInnovation(ans)

		n.innovation++

		return ans
	}

	nodeMutate := func() {
		var firstNode int
		var secondNode int
		ans := false

		//todo find a better way to check (for both statements)
		for !ans {
			firstNode = int(r.Int63n(int64(nodeRange+1)))
			secondNode = int(r.Int63n(int64(nodeRange+1)))

			for i := 0; i < len(n.connectionInnovation); i++ {
				if n.connectionInnovation[i][0] == firstNode && n.connectionInnovation[i][1] == secondNode {
					ans = true
				}
			}
		}

		network.mutateNode(firstNode, secondNode)

		addConnectionInnovation(firstNode, secondNode)
	}
	//todo fix the casting
	if r.Float64() <= n.nodeMutate {
		nodeMutate()
	} else {
		var firstNode int
		var secondNode int
		ans := false
		attempts := 0
		for !ans && attempts <= 5 {
			firstNode = int(r.Int63n(int64(nodeRange+1)))
			secondNode = int(r.Int63n(int64(nodeRange+1)))

			ans = true
			for i := 0; i < len(n.connectionInnovation); i++ {
				if n.connectionInnovation[i][0] == firstNode && n.connectionInnovation[i][1] == secondNode {
					ans = false
				}
			}

			attempts++
		}

		if(attempts > 5) {
			nodeMutate()
		}

		addConnectionInnovation(firstNode, secondNode)

		network.mutateConnection(int(r.Int63n(int64(nodeRange+1))), int(r.Int63n(int64(nodeRange+1))))
	}
}
