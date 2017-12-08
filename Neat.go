package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

//MAX 1000 innovation
/*
not going to speciate until after a couple of rounds
*/

type Neat struct {
	connectMutate        float64   //odds for connection mutation
	nodeMutate           float64   //odds for node mutation
	innovation           int       //number of innovations
	network              []Network //stores networks in species
	connectionInnovation [][]int   //stores innovation number and connection to and from ex: 1, fromNode:2, toNode: 5 [2,5]
	speciesThreshold     float64   //could adjust based upon average difference between networks
	networkId            int
	species              []Species
	numSpecies           int
	speciesId            int
}

func GetNeatInstance(numNetworks int, input int, output int) Neat {
	n := Neat{innovation: 0, connectMutate: .7, speciesThreshold: .01,
		nodeMutate: .3, network: make([]Network, numNetworks), connectionInnovation: make([][]int, 0, 1000), species: make([]Species, 0, 5)}

	//TODO: make sure correct
	for i := output; i < input+output; i++ {
		for a := 0; a < output; a++ {
			n.createNewInnovation([]int{i, a})
		}
	}

	for i := 0; i < len(n.network); i++ {
		n.network[i] = GetNetworkInstance(input, output, i, 0, .1)
	}

	n.createSpecies(n.network[0 : len(n.network)%5+(numNetworks/5)])
	for i := len(n.network)%5 + (numNetworks / 5); i+(numNetworks/5) <= len(n.network); i += (numNetworks / 5) {
		n.createSpecies(n.network[i : i+(numNetworks/5)])
	}

	n.mutatePopulation()
	n.speciateAll()
	n.checkSpecies()

	return n
}

func (n *Neat) speciateAll() {
	for i := 0; i < len(n.network); i++ {
		n.speciate(&n.network[i])
	}
}

func (n *Neat) speciate(network *Network) {
	values := make([]float64, len(n.species))

	for i := 0; i < len(n.species); i++ {
		values[i] = compareGenome(network.id+1, network.innovation, n.species[i].commonNodes, n.species[i].commonConnection)
	}

	//this should be faster than sorting the whole thing (it also retains position information)
	index := 0
	lValue := 1000.0
	for i := 0; i < len(values); i++ {
		if values[i] < lValue {
			index = i
			lValue = values[i]
		}
	}

	if lValue > n.speciesThreshold { //i flipped this sign i think it works better %different > differentThreshold
		//finds the position
		networkIndex := 0
		for i := 0; i < len(n.network); i++ {
			if n.network[i].networkId == network.networkId {
				networkIndex = i
			}
		}

		specId := network.species

		newSpec := n.createSpecies(n.network[networkIndex : networkIndex+1]) //TODO: need to remove network from species?

		s := n.getSpecies(specId)
		if s != nil {
			s.removeNetwork(network.networkId)
			for i := 0; i < len(s.network); i++ {
				if s.network[i].networkId != network.networkId {
					n.speciate(s.network[i])
				}
			}

			if len(s.network) < 2 {
				n.removeSpecies(s.id)
			}
		}

		if len(newSpec.network) < 2 {
			newSpec.removeNetwork(network.networkId)

			n.species[index].addNetwork(network)

			n.removeSpecies(newSpec.id)
		}
	} else {
		spec := n.getSpecies(network.species)
		if spec != nil {
			spec.removeNetwork(network.networkId)
		}
		n.species[index].addNetwork(network)
	}
}

func (n *Neat) getSpecies(id int) *Species {
	for i := 0; i < len(n.species); i++ {
		if isRealSpecies(&n.species[i]) && n.species[i].id == id {
			return &n.species[i]
		}
	}

	return nil
}

func compareGenome(node int, innovation []int, nodeA int, innovationA []int) float64 {
	var larger []int
	var smaller []int

	if len(innovation) > len(innovationA) {
		larger = innovation
		smaller = innovationA
	} else {
		larger = innovationA
		smaller = innovation
	}

	missing := 0
	for b := 0; b < len(larger); b++ {
		if sort.SearchInts(smaller, larger[b]) == len(smaller) {
			missing++
		}
	}

	return float64(missing+int(math.Abs(float64(node-nodeA)))) / (float64(len(smaller)) + float64((node+nodeA)/2))
}

func (n *Neat) findInnovationNum(search []int) int {
	for i := 0; i < len(n.connectionInnovation); i++ {
		if n.connectionInnovation[i][0] == search[0] && n.connectionInnovation[i][1] == search[1] {
			return i
		}
	}

	return -1
}

func (n *Neat) checkSpecies() {
	for i := 0; i < len(n.species); i++ {
		values := make([]float64, len(n.species))
		for a := 0; a < len(n.species); a++ {
			if a == i {
				values[a] = 100.0
				continue
			}

			values[a] = compareGenome(n.species[i].commonNodes, n.species[i].commonConnection, n.species[a].commonNodes, n.species[a].commonConnection)
		}

		lValue := 1000.0
		for a := 0; a < len(values); a++ {
			if values[a] < lValue {
				lValue = values[a]
			}
		}
		if lValue < n.speciesThreshold || n.species[i].numNetwork < 2 { //switched direction if sign because %dif < difthreshold for it to be the same
			n.removeSpecies(n.species[i].id)
			/*currentSpecies := n.species[i].network
			n.species = append(n.species[:i], n.species[(i+1):]...)
			for a := 0; a < len(currentSpecies); a++ {
				n.speciate(currentSpecies[a])
			}

			continue*/
		}
	}
}

func (n *Neat) mutatePopulation() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numNet := int(r.Int63n(int64(len(n.network)))-3/5) + 3
	for i := 0; i < numNet; i++ {
		species := int(r.Int63n(int64(len(n.species))))
		network := n.species[species].getNetworkAt(int(r.Int63n(int64(n.species[species].numNetwork))))

		nodeRange := network.id

		addConnectionInnovation := func(numFrom int, numTo int) int {
			//checks to see if preexisting innovation
			for i := 0; i < len(n.connectionInnovation); i++ {
				if n.connectionInnovation[i][1] == numTo && n.connectionInnovation[i][0] == numFrom {
					//network.addInnovation(i)
					n.species[species].mutateNetwork(i)

					return i
				}
			}

			//checks to see if needs to grow
			num := n.createNewInnovation([]int{numFrom, numTo})

			//network.addInnovation(num)
			n.species[species].mutateNetwork(num)

			return num
		}

		nodeMutate := func() {
			var firstNode int
			var secondNode int
			ans := false

			for !ans {
				firstNode = int(rand.Float64() * float64(nodeRange))

				if !isOutput(network.getNode(firstNode)) {
					ans = true
				}
			}

			secondNode = network.getNode(firstNode).send[int(rand.Float64()*float64(len(network.getNode(firstNode).send)))].nodeTo.id //int(r.Int63n(int64(nodeRange)))

			a := addConnectionInnovation(firstNode, network.getNextNodeId())
			b := addConnectionInnovation(network.getNextNodeId(), secondNode)
			network.mutateNode(firstNode, secondNode, a, b)
			n.species[species].nodeCount++
		}

		if r.Float64() <= n.nodeMutate {
			nodeMutate()
		} else {
			var firstNode int
			var secondNode int
			ans := true
			attempts := 0
			for ans && attempts <= 10 {
				firstNode = int(r.Int63n(int64(nodeRange)))
				secondNode = int(r.Int63n(int64(nodeRange)))

				if firstNode == secondNode || isOutput(network.getNode(firstNode)) || isInput(network.getNode(secondNode)) {
					continue
				}

				ans = true
				for i := 0; i < len(n.connectionInnovation); i++ {
					if n.connectionInnovation[i][0] == firstNode && n.connectionInnovation[i][1] == secondNode || n.connectionInnovation[i][1] == firstNode && n.connectionInnovation[i][0] == secondNode {

						ans = network.containsInnovation(i)
					}
				}

				attempts++
			}

			if attempts > 10 {
				nodeMutate()
				continue
			}

			network.mutateConnection(firstNode, secondNode, addConnectionInnovation(firstNode, secondNode))
		}
	}
}

//TODO: need to have looped until completion
func (n *Neat) start(input [][][]float64) {
	for i := 0; i < len(n.species); i++ {
		if isRealSpecies(&n.species[i]) {
			n.species[i].trainNetworks(input)
		}
	}

	n.checkSpecies()

	newOveral := make([]Network, len(n.network))
	count := 0
	for i := 0; i < len(n.species); i++ {
		if isRealSpecies(&n.species[i]) {
			newNets := n.species[i].mateSpecies()
			for a := 0; a < len(newNets); a++ {
				newOveral[count] = newNets[a]
				count++
			}
		}
	}
}

func (n *Neat) printNeat() {
	fmt.Println()
	fmt.Println()
	for i := 0; i < len(n.species); i++ {
		fmt.Println("species id: ", n.species[i].id)
		for a := 0; a < len(n.species[i].network); a++ {
			fmt.Println("network id: ", n.species[i].network[a].networkId, " species id: ", n.species[i].network[a].species)

			fmt.Print("expected connection: ", n.species[i].network[a].innovation)
			for b := 0; b < len(n.species[i].network[a].innovation); b++ {
				fmt.Print(n.getInnovation(n.species[i].network[a].innovation[b]))
			}
			fmt.Println()

			for b := 0; b < len(n.species[i].network[a].nodeList); b++ {
				fmt.Print("node: ", n.species[i].network[a].nodeList[b].id, " sending: ")
				for c := 0; c < len(n.species[i].network[a].nodeList[b].send); c++ {
					fmt.Print(n.species[i].network[a].nodeList[b].send[c].nodeTo.id, " ")
				}

				fmt.Print("receive: ")
				for c := 0; c < len(n.species[i].network[a].nodeList[b].receive); c++ {
					fmt.Print(n.species[i].network[a].nodeList[b].receive[c].nodeFrom.id, " ")
				}

				fmt.Println()
			}
		}
	}
}

func (n *Neat) createNewInnovation(values []int) int {
	if n.innovation > cap(n.connectionInnovation)-1 {
		n.connectionInnovation = append(n.connectionInnovation, values)
	} else {
		n.connectionInnovation = n.connectionInnovation[0 : len(n.connectionInnovation)+1]
		n.connectionInnovation[n.innovation] = values
	}
	n.innovation++

	return n.innovation - 1
}

func (n *Neat) getInnovation(num int) []int {
	return n.connectionInnovation[num]
}

func (n *Neat) createSpecies(possible []Network) *Species {
	for i := 0; i < len(possible); i++ {
		possible[i].species = n.speciesId
	}

	s := GetSpeciesInstance(n.speciesId, possible, &n.connectionInnovation)
	if cap(n.species) <= len(n.species) {
		n.species = append(n.species, s)
	} else {
		n.species = n.species[0 : len(n.species)+1]
		n.species[len(n.species)-1] = s
	}

	n.numSpecies++
	n.speciesId++

	return &n.species[len(n.species)-1]
}

func (n *Neat) removeSpecies(id int) {
	for i := 0; i < len(n.species); i++ {
		if n.species[i].id == id {
			currentSpecies := n.species[i].network
			n.species = append(n.species[:i], n.species[i+1:]...)
			for a := 0; a < len(currentSpecies); a++ {
				n.speciate(currentSpecies[a])
			}
		}
	}
}
