package main

import (
	"math/rand"
	"sort"
	"time"
	//	"fmt"
	"sync"
)

type Species struct {
	network             []*Network //holds the pointer to all the networks
	connectionInnovaton []int      //holds number of occerences of each innovation
	commonInnovation    []int      //common connection innovation numbers
	innovationDict      *[][]int   //master list for all innovations
	id                  int        //the identifier for the species
	mutate              float64
}

func GetSpeciesInstance(id int, networks []Network, innovations *[][]int, mutate float64) Species {
	s := Species{mutate: mutate, id: id, network: make([]*Network, len(networks)), commonInnovation: make([]int, 0, len(*innovations)*2), connectionInnovaton: make([]int, len(*innovations)*2), innovationDict: innovations}

	for i := 0; i < len(networks); i++ {
		s.network[i] = &networks[i]
	}

	s.updateStereotype()

	return s
}

////////////////////////////////////////////////////////////INNOVATION
func (s *Species) addCI(a int) {
	for i := 0; i < len(s.commonInnovation); i++ {
		if s.commonInnovation[i] == a {
			return
		}
	}

	if len(s.commonInnovation) >= cap(s.commonInnovation) {
		s.commonInnovation = append(s.commonInnovation, a)
	} else {
		s.commonInnovation = s.commonInnovation[0: len(s.commonInnovation)+1]
		s.commonInnovation[len(s.commonInnovation)-1] = a
	}
}
func (s *Species) removeCI(a int) {
	for i := 0; i < len(s.commonInnovation); i++ {
		if s.commonInnovation[i] == a {
			s.commonInnovation = append(s.commonInnovation[:i], s.commonInnovation[i+1:]...)
		}
	}
}
func (s *Species) getInovOcc(i int) *int {
	if i >= len(s.connectionInnovaton) {
		insert := make([]int, 1+i-len(s.connectionInnovaton))
		s.connectionInnovaton = append(s.connectionInnovaton, insert...)
	}
	return &s.connectionInnovaton[i]
}
func (s *Species) incrementInov(i int) *int {
	ans := s.getInovOcc(i)
	*ans++

	if float64(*ans)/float64(len(s.network)) >= .5 { //could have issues
		//if float64(s.connectionInnovaton[i]/len(s.network)) > .6 {
		s.addCI(i)
	}

	return ans
}
func (s *Species) reduceInov(i int) *int {
	ans := s.getInovOcc(i)
	*ans--

	if float64(*ans)/float64(len(s.network)) < .5 { //could have issues
		s.removeCI(i)
	}

	return ans
}
func (s *Species) checkCI() {
	for i := 0; i < len(s.commonInnovation); i++ {
		s.removeCI(s.commonInnovation[i])
	}

	for i := 0; i < len(s.connectionInnovaton); i++ {
		if float64(s.connectionInnovaton[i])/float64(len(s.network)) >= .5 {
			s.addCI(i)
		}
	}
}
func (n *Species) getInnovationRef(num int) []int {
	return (*n.innovationDict)[num]
}
func (s *Species) sortInnovation() {
	for i := 0; i < len(s.network); i++ {
		sort.Ints(s.network[i].innovation)
	}
}

//////////////////////////////////////////////////////////////NETWORK
func (s *Species) getNetworkAt(a int) *Network {
	return s.network[a]
}
func (s *Species) removeNetwork(id int) {
	for i := 0; i < len(s.network); i++ {
		if s.network[i].networkId == id {
			inn := s.network[i].innovation
			s.network = append(s.network[:i], s.network[i+1:]...)

			for a := 0; a < len(inn); a++ {
				s.reduceInov(inn[a])
			}

			s.checkCI()
		}
	}
}
func (s *Species) getNetwork(id int) *Network {
	for i := 0; i < len(s.network); i++ {
		if s.network[i] != nil && s.network[i].networkId == id {
			return s.network[i]
		}
	}

	return nil
}
func (s *Species) addNetwork(n *Network) {
	if len(s.network) >= cap(s.network) {
		s.network = append(s.network, n)
	} else {
		s.network = s.network[0: len(s.network)+1]
		s.network[len(s.network)-1] = n
	}

	n.species = s.id

	for i := 0; i < len(n.innovation); i++ {
		s.incrementInov(n.innovation[i])
	}

	s.checkCI()
}

///////////////////////////////////////////////////////////MATE+MUTATE
func (s *Species) updateStereotype() {
	numNodes := 0

	for i := 0; i < len(s.connectionInnovaton); i++ {
		s.connectionInnovaton[i] = 0
	}

	for i := 0; i < len(s.commonInnovation); i++ {
		s.removeCI(s.commonInnovation[i])
	}

	for i := 0; i < len(s.network); i++ {
		if s.network[i] != nil {
			numNodes += len(s.network[i].nodeList)
			for a := 0; a < len(s.network[i].innovation); a++ {
				s.incrementInov(s.network[i].innovation[a])
			}
		}
	}
}
func (n *Species) createNewInnovation(values []int) int {
	*n.innovationDict = (*n.innovationDict)[0: len(*n.innovationDict)+1]
	(*n.innovationDict)[len(*n.innovationDict)-1] = values

	return len(*n.innovationDict) - 1
}
func (s *Species) mutateNetwork(network *Network, nodeMutateA float64) {
	nodeRange := len(network.nodeList)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	addConnectionInnovation := func(numFrom int, numTo int) int {
		//checks to see if preexisting innovation
		for i := 0; i < len((*s.innovationDict)); i++ {
			if (*s.innovationDict)[i][1] == numTo && (*s.innovationDict)[i][0] == numFrom {
				//network.addInnovation(i)
				s.incrementInov(i)

				return i
			}
		}

		//checks to see if needs to grow
		num := s.createNewInnovation([]int{numFrom, numTo})

		//network.addInnovation(num)
		s.incrementInov(num)

		return num
	}

	nodeMutate := func() {
		var firstNode int
		var secondNode int
		ans := false

		for !ans {
			firstNode = int(rand.Float64() * float64(nodeRange))

			if !isOutput(network.getNode(firstNode)) && len(network.getNode(firstNode).send) > 0 {
				ans = true
			}
		}

		secondNode = network.getNode(firstNode).send[int(rand.Float64()*float64(len(network.getNode(firstNode).send)))].nodeTo.id //int(r.Int63n(int64(nodeRange)))

		network.mutateNode(firstNode, secondNode, addConnectionInnovation(firstNode, network.getNextNodeId()), addConnectionInnovation(network.getNextNodeId(), secondNode))
	}

	if r.Float64() <= nodeMutateA {
		nodeMutate()
	} else {
		/*
			could interate through and find a number that has not been used and then use that number so only have to rng one
		*/
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

			ans = false
			for i := 0; i < len((*s.innovationDict)); i++ {
				if ((*s.innovationDict)[i][0] == firstNode && (*s.innovationDict)[i][1] == secondNode) || ((*s.innovationDict)[i][1] == firstNode && (*s.innovationDict)[i][0] == secondNode) {
					ans = network.containsInnovation(i)
					if ans {
						break
					}
				}
			}

			if !ans {
				ans = network.checkCircleMaster(network.getNode(firstNode), secondNode)
			}

			attempts++
		}

		if attempts > 10 {
			nodeMutate()
		} else {
			network.mutateConnection(firstNode, secondNode, addConnectionInnovation(firstNode, secondNode))
		}
	}
}
func (n *Species) mateNetwork(nB Network, nA Network) Network {
	ans := GetNetworkInstance(len(nB.output), len(nB.input)-1, 0, nB.species, nB.learningRate, false)

	numNode := -1 * (len(nB.output) + len(nB.input))
	if len(nA.nodeList) > len(nB.nodeList) {
		numNode += len(nA.nodeList)
	} else {
		numNode += len(nB.nodeList)
	}

	for i := 0; i < numNode; i++ { //this should be ok
		ans.createNode(100)
	}

	for i := 0; i < len(nA.innovation); i++ {
		ans.mutateConnection(n.getInnovationRef(nA.getInovation(i))[0], n.getInnovationRef(nA.getInovation(i))[1], nA.getInovation(i))
	}

	for i := 0; i < len(nB.innovation); i++ {
		if !ans.containsInnovation(nB.innovation[i]) {
			ans.mutateConnection(n.getInnovationRef(nB.getInovation(i))[0], n.getInnovationRef(nB.getInovation(i))[1], nB.getInovation(i))
		}
	}

	return ans
}
func (s *Species) trainNetworks(trainingSet [][][]float64, control *sync.WaitGroup) {
	for i := 0; i < len(s.network); i++ {
		s.network[i].trainSet(trainingSet, 1000)
	}
	control.Done()
}
//used to make networks inside a species
func (s *Species) mateSpecies() []Network {
	s.adjustFitness()

	//sorts by adjusted fitness
	sortedNetwork := make([]*Network, len(s.network)*85/100)
	lastValue := 1000.0
	sumFitness := 0.0
	for i := 0; i < len(sortedNetwork); i++ {
		localMax := 0.0
		localIndex := 0
		for a := 0; a < len(s.network); a++ {
			if s.getNetworkAt(a).adjustedFitness > localMax && s.getNetworkAt(a).adjustedFitness <= lastValue {
				good := true
				for b := i - 1; b >= 0; b-- {
					if s.getNetworkAt(a).networkId == sortedNetwork[b].networkId {
						good = false
						break
					}

					if sortedNetwork[b].adjustedFitness != s.getNetworkAt(a).adjustedFitness {
						break
					}
				}

				if good {
					localMax = s.network[a].adjustedFitness
					localIndex = a
				}
			}
		}

		sortedNetwork[i] = s.getNetworkAt(localIndex)
		sumFitness += sortedNetwork[i].adjustedFitness
		lastValue = sortedNetwork[i].adjustedFitness
	}

	newNets := make([]Network, len(s.network))
	count := 0
	for i := 0; i < len(sortedNetwork); i++ {
		numKids := int(sortedNetwork[i].adjustedFitness / sumFitness * float64(len(newNets)))
		numMade := numKids
		for a := 1; count < len(newNets) && a+i < len(sortedNetwork); a++ {
			if sortedNetwork[i+a] != nil {
				newNets[count] = s.mateNetwork(*sortedNetwork[i], *sortedNetwork[i+a])
				count++
				numMade--
			}
		}

	}

	for i := 0; count < len(newNets); i++ {
		s.mutateNetwork(sortedNetwork[i], s.mutate) //adds best network back in where the last child for that network
		sortedNetwork[i].resetWeight()
		newNets[count] = *sortedNetwork[i]
		count++

		if i == len(sortedNetwork)-1 {
			i-- //this can lead to mutating the same network as last time (stacking mutations) but i don't think it is a big deal
		}
	}

	for i := 0; i < len(newNets); i++ {
		newNets[i].networkId = s.network[i].networkId
	}

	s.updateStereotype()

	return newNets
}
func (s *Species) adjustFitness() {
	for i := 0; i < len(s.network); i++ {
		s.network[i].adjustedFitness = s.network[i].fitness / float64(len(s.network))
	}
}
func (s *Species) avgNode() int {
	if len(s.network) == 0 {
		return 0
	}
	sum := 0
	for i := 0; i < len(s.network); i++ {
		sum += len(s.network[i].nodeList)
	}

	return sum / len(s.network)
}