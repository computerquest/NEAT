package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

/*
not going to speciate until after a couple of rounds
*/

type Neat struct {
	connectMutate        float64   //odds for connection mutation
	nodeMutate           float64   //odds for node mutation
	innovation           int       //number of innovations
	network              []Network //stores networks in species
	connectionInnovation [][]int   //stores innovation number and connection to and from ex: 1, fromNode:2, toNode: 5
	speciesThreshold     float64   //could adjust based upon average difference between networks
	networkId            int
	species              []Species
	numSpecies           int
}

//TODO: fix id system
//TODO: finish
func GetNeatInstance(numNetworks int, input int, output int) Neat {
	n := Neat{innovation: 0, connectMutate: .7,
		nodeMutate: .3, network: make([]Network, numNetworks), connectionInnovation: make([][]int, 10), species: make([]Species, 5)}

	for i := 0; i < len(n.connectionInnovation); i++ {
		n.connectionInnovation[i] = make([]int, 2)
	}

	/*REST OF METHOD
	create species
	perform initial mutations
		between 1-3 for each network
	speciate
	*/

	for i := 0; i < len(n.network); i++ {
		n.network[len(n.network)-1-i] = GetNetworkInstance(input, output, i, 0)
	}

	n.species[0] = GetSpeciesInstance(100, n.network[0:len(n.network)%5+(numNetworks/5)+1])
	for i, b := len(n.network)%5+(numNetworks/5)+1, 1; i+(numNetworks/5) < len(n.network); i, b = i+(numNetworks/5), b+1 {
		n.species[b] = GetSpeciesInstance(100, n.network[i:i+(numNetworks/5)])
		//TODO: uncomment when done
		//n.mutateNetwork()
	}

	return n
}

func (n *Neat) mateNetwork(nB Network, nA Network, idNum int) Network {
	ans := GetNetworkInstance(len(nB.output), len(nB.input), idNum, nB.species)

	var numNode int
	if nA.id > nB.id {
		numNode = nA.id
	} else {
		numNode = nB.id
	}

	for i := ans.id; i < numNode; i++ {
		ans.createNode()
	}

	for i := 0; i < nA.numInnovation; i++ {
		ans.mutateConnection(n.getInnovation(nA.getInovation(i))[0], n.getInnovation(nA.getInovation(i))[1], nA.getInovation(i))
	}

	for i := 0; i < nB.numInnovation; i++ {
		exist := false
		for a := 0; a < nA.numInnovation; a++ {
			if nB.getInovation(i) == nA.getInovation(a) {
				exist = true
				break
			}
		}

		if !exist {
			ans.mutateConnection(n.getInnovation(nB.getInovation(i))[0], n.getInnovation(nB.getInovation(i))[1], nB.getInovation(i))
		}
	}

	return ans
}
func (n *Neat) mateSpecies(s *Species) {
	s.adjustFitness()

	//sorts by adjusted fitness
	sortedNetwork := make([]*Network, s.numNetwork*85/100)
	sumFitness := 0.0
	for i := 0; i < len(s.network); i++ {
		for a := i; a < len(s.network); a++ {
			if s.network[i] != nil && s.getNetworkAt(a).adjustedFitness > s.network[i].adjustedFitness {
				sortedNetwork[i] = s.getNetworkAt(a)
			}
		}

		sumFitness += sortedNetwork[i].adjustedFitness
	}

	for i := 1; i < len(sortedNetwork); i++ {

	}
}

//rewrite
/*
are you comparing every network to every other or are you comparing random geneomes (collection of genes) from last generation of species to each network
*/
/*
could improve the results by:
creating a custom rep network that generalizes the species (a stereotype if you will)
pay greater attention to the speciation comparisons and compare values later in case of dispute
*/
//TODO: need to make sure that even the connections made when initialized are included in the innovation numbers here
//TODO: need plan for creating a new species and rearranging species
//TODO: need a plan for starting a new species
//TODO: finish
func (n *Neat) speciate(network *Network) {
	/*repNetworks := make([]*Network, n.species)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < len(n.network); i++ {
		repNetworks[i] = &n.network[i][int(r.Int63n(int64(n.nps)))]
	}

	newNet := make([][]Network, n.species)
	for i := 0; i < len(n.network); i++ {
		newNet[i] = make([]Network, n.nps)
	}

	calcSpecies := func(inputS []int, inputL []int) int {
		missing := 0
		for b := 0; b < len(inputS); b++ {
			ans := sort.SearchInts(inputL, inputS[b])

			//TODO: find default return value
			if ans == -1 {
				missing++
			}
		}

		return missing
	}

	//TODO: neaten up
	for i := 0; i < len(n.network); i++ {
		for a := 0; a < len(n.network[i]); a++ {
			values := make([]float64, n.species)

			for z := 0; z < len(repNetworks); z++ {
				if len(repNetworks[z].innovation) < len(n.network[i][a].innovation) {
					values[z] = float64(calcSpecies(repNetworks[z].innovation, n.network[i][a].innovation) / len(n.network[i][a].innovation))
				} else {
					values[z] = float64(calcSpecies(n.network[i][a].innovation, repNetworks[z].innovation) / len(repNetworks[z].innovation))
				}

			}
			min := 2.0
			index := 0
			for b := 0; b < len(values); b++ {
				if values[b] < min {
					index = b
					min = values[b]
				}
			}

			newNet[index][len(newNet[index])-1] = n.network[i][a]
		}
	}*/
	values := make([]float64, len(n.species))

	for i := 0; i < len(n.species); i++ {
		if len(n.species[i].connectionInnovaton) > len(network.innovation) {
			values[i] = compareGenome(network.id+1, network.innovation, n.species[i].commonNodes, n.species[i].commonConnection)
		} else {
			values[i] = compareGenome(n.species[i].commonNodes, n.species[i].commonConnection, network.id+1, network.innovation)
		}
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

	if lValue < n.speciesThreshold {
		fmt.Print(index)
	} else {

	}
}

//TODO: STUFF
//recieves input in order shortest to longest
func compareGenome(node int, innovation []int, nodeA int, innovationA []int) float64 {
	missing := 0
	for b := 0; b < len(innovation); b++ {
		ans := sort.SearchInts(innovationA, innovation[b])

		//TODO: verify default return value
		if ans == len(innovationA) {
			missing++
		}
	}

	return float64((missing + int(math.Abs(float64(node-nodeA)))) / (len(innovationA) + int((node+nodeA)/2)))
}

//TODO: finalize protocol for same species
func (n *Neat) checkSpecies() {
	for i := 0; i < len(n.species); i++ {
		values := make([]float64, len(n.species))
		for a := 0; a < len(n.species); a++ {
			if a == i {
				continue
			}

			if len(n.species[i].commonConnection) < len(n.species[a].commonConnection) {
				values[a] = compareGenome(n.species[i].commonNodes, n.species[i].commonConnection, n.species[a].commonNodes, n.species[a].commonConnection)
			} else {
				values[a] = compareGenome(n.species[a].commonNodes, n.species[a].commonConnection, n.species[i].commonNodes, n.species[i].commonConnection)
			}
		}

		index := 0
		lValue := 1000.0
		for i := 0; i < len(values); i++ {
			if values[i] < lValue {
				index = i
				lValue = values[i]
			}
		}
		fmt.Print(index)
		if lValue > n.speciesThreshold {
			currentSpecies := n.species[i].network
			n.species = append(n.species[:i], n.species[(i+1):]...)
			for a := 0; a < len(currentSpecies); a++ {
				n.speciate(currentSpecies[a])
			}
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

		//TODO: test
		addConnectionInnovation := func(numTo int, numFrom int) int {
			ans := n.innovation
			if len(n.connectionInnovation) <= (n.innovation + 1) {
				newStuff := []int{numFrom, numTo}
				n.connectionInnovation = append(n.connectionInnovation, newStuff)
			} else {
				n.connectionInnovation[n.innovation][0] = numFrom
				n.connectionInnovation[n.innovation][1] = numTo
			}

			network.addInnovation(ans)

			n.species[species].mutateNetwork(ans)

			n.innovation++

			return ans
		}

		nodeMutate := func() {
			var firstNode int
			var secondNode int
			ans := false

			//TODO: find a better way to check (for both statements)
			for !ans {
				firstNode = int(r.Int63n(int64(nodeRange + 1)))
				secondNode = int(r.Int63n(int64(nodeRange + 1)))

				for i := 0; i < len(n.connectionInnovation); i++ {
					if n.connectionInnovation[i][0] == firstNode && n.connectionInnovation[i][1] == secondNode {
						ans = true
					}
				}
			}

			//TODO: give actual innovation numbers
			network.mutateNode(firstNode, secondNode, 100, 100)
			n.species[species].nodeCount++

			addConnectionInnovation(firstNode, secondNode)
		}
		//TODO: fix the casting
		if r.Float64() <= n.nodeMutate {
			nodeMutate()
		} else {
			var firstNode int
			var secondNode int
			ans := false
			attempts := 0
			for !ans && attempts <= 5 {
				firstNode = int(r.Int63n(int64(nodeRange + 1)))
				secondNode = int(r.Int63n(int64(nodeRange + 1)))

				ans = true
				for i := 0; i < len(n.connectionInnovation); i++ {
					if n.connectionInnovation[i][0] == firstNode && n.connectionInnovation[i][1] == secondNode {
						ans = false
					}
				}

				attempts++
			}

			if attempts > 5 {
				nodeMutate()
			}

			addConnectionInnovation(firstNode, secondNode)

			//TODO: change the connection number
			network.mutateConnection(int(r.Int63n(int64(nodeRange+1))), int(r.Int63n(int64(nodeRange+1))), 100)
		}
	}
}

//TODO: make sure all changes have been made to real method
func (n *Neat) mutatePopulationTest() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numNet := int(r.Int63n(int64(len(n.network)))-3/5) + 3
	for i := 0; i < numNet; i++ {
		species := 0
		network := n.species[species].getNetworkAt(0)

		nodeRange := network.id

		//TODO: test
		addConnectionInnovation := func(numTo int, numFrom int) int {
			ans := n.innovation
			if len(n.connectionInnovation) <= (n.innovation + 1) {
				newStuff := []int{numFrom, numTo}
				n.connectionInnovation = append(n.connectionInnovation, newStuff)
			} else {
				n.connectionInnovation[n.innovation][0] = numFrom
				n.connectionInnovation[n.innovation][1] = numTo
			}

			network.addInnovation(ans)

			n.species[species].mutateNetwork(ans)

			n.innovation++

			return ans
		}

		nodeMutate := func() {
			var firstNode int
			var secondNode int
			ans := false

			//TODO: find a better way to check (for both statements)
			for !ans {
				firstNode = int(r.Int63n(int64(nodeRange + 1)))
				secondNode = int(r.Int63n(int64(nodeRange + 1)))

				for i := 0; i < len(n.connectionInnovation); i++ {
					if n.connectionInnovation[i][0] == firstNode && n.connectionInnovation[i][1] == secondNode {
						ans = true
					}
				}
			}

			//TODO: give actual innovation numbers
			network.mutateNode(firstNode, secondNode, 25, 25)
			n.species[species].nodeCount++

			addConnectionInnovation(firstNode, secondNode)
		}
		//TODO: fix the casting
		if r.Float64() <= n.nodeMutate {
			nodeMutate()
		} else {
			var firstNode int
			var secondNode int
			ans := false
			attempts := 0
			for !ans && attempts <= 5 {
				firstNode = int(r.Int63n(int64(nodeRange + 1)))
				secondNode = int(r.Int63n(int64(nodeRange + 1)))

				ans = true
				for i := 0; i < len(n.connectionInnovation); i++ {
					if n.connectionInnovation[i][0] == firstNode && n.connectionInnovation[i][1] == secondNode {
						ans = false
					}
				}

				attempts++
			}

			if attempts > 5 {
				nodeMutate()
			}

			addConnectionInnovation(firstNode, secondNode)

			//TODO: change the connection number
			network.mutateConnection(int(r.Int63n(int64(nodeRange+1))), int(r.Int63n(int64(nodeRange+1))), 25)
		}
	}
}

func (n *Neat) createNewInnovation(values []int) []int {
	n.innovation++
	if n.innovation > len(n.connectionInnovation)-1 {
		n.connectionInnovation = append(n.connectionInnovation, values)
	}

	return n.getInnovation(n.innovation)
}

func (n *Neat) getInnovation(num int) []int {
	return n.connectionInnovation[len(n.connectionInnovation)-1-num]
}
