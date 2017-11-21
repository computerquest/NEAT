package main

import ()

/*
not going to speciate until after a couple of rounds
*/

//TODO: more robust species id system.
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
	speciesId            int
}

//TODO: fix id system ?
func GetNeatInstance(numNetworks int, input int, output int) Neat {
	n := Neat{innovation: 0, connectMutate: .7,
		nodeMutate: .3, network: make([]Network, numNetworks), connectionInnovation: make([][]int, 10), species: make([]Species, 5)}

	for i := 0; i < len(n.connectionInnovation); i++ {
		n.connectionInnovation[i] = make([]int, 2)
	}

	for i := 0; i < len(n.network); i++ {
		n.network[len(n.network)-1-i] = GetNetworkInstance(input, output, i, 0)
	}

	n.createSpecies(n.network[0 : len(n.network)%5+(numNetworks/5)+1])
	for i, b := len(n.network)%5+(numNetworks/5)+1, 1; i+(numNetworks/5) < len(n.network); i, b = i+(numNetworks/5), b+1 {
		n.createSpecies(n.network[i : i+(numNetworks/5)])
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

//TODO: why do i have this method
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

//TODO: need a plan for starting a new species
//TODO: finish/test
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

	//TODO: create new species
	if lValue < n.speciesThreshold {
		fmt.Print(index)
	} else {
		//TODO: need more accurate way to change species (search by id)
		n.species[network.species] = nil
		n.species[index].addNetwork(network)
	}
}

//recieves input in order shortest to longest
func compareGenome(node int, innovation []int, nodeA int, innovationA []int) float64 {
	var larger []int
	var smaller []int

	if len(innovation) > innovationA {
		larger = innovation
		smaller = innovationA
	} else {
		larger = innovationA
		smaller = innovation
	}

	missing := 0
	for b := 0; b < len(larger); b++ {
		ans := sort.SearchInts(smaller, larger[b])

		//TODO: verify default return value
		if ans == len(smaller) {
			missing++
		}
	}

	return float64((missing + int(math.Abs(float64(node-nodeA)))) / (len(smaller) + int((node+nodeA)/2)))
}

//TODO: test
func (n *Neat) checkSpecies() {
	for i := 0; i < len(n.species); i++ {
		values := make([]float64, len(n.species))
		for a := 0; a < len(n.species); a++ {
			if a == i {
				continue
			}

			values[a] = compareGenome(n.species[i].commonNodes, n.species[i].commonConnection, n.species[a].commonNodes, n.species[a].commonConnection)
		}

		index := 0
		lValue := 1000.0
		for i := 0; i < len(values); i++ {
			if values[i] < lValue {
				index = i
				lValue = values[i]
			}
		}
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

//TODO: test
func (n *Neat) createSpecies(possible []*Network) {
	s := GetSpeciesInstance(n.speciesId, possible, n.innovation)
	if n.numSpecies >= len(n.species) {
		n.species = append(n.species, s)
	} else {
		n[len(n.species)-n.numSpecies-1] = s
	}

	n.numSpecies++
	n.speciesId++
}

//TODO: test
func (n *Neat) removeSpecies(id int) {
	for i := 0; i < len(n.species); i++ {
		if n.species[i].id == id {
			currentSpecies := n.species[i].network
			s.network = append(s.network[:i], s.network[i+1:]...)
			for a := 0; a < len(currentSpecies); a++ {
				n.speciate(currentSpecies[a])
			}
		}
	}
}
