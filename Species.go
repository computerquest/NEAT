package main

import (
	"sort"
)

//TODO: might want to consider starting the innovation master list at one so that all of the arrays have a default value (or prevents default value)
//TODO: look into node id system and make sure that it doesn't allow different types of nodes to have the same id (screw innovation number pairings)
//TODO: fix the avg because empty slots created by append will screw
type Species struct {
	network             []*Network //holds the pointer to all the networks
	connectionInnovaton []int      //holds number of occerences of each innovation
	nodeCount           int        //number of total nodes
	commonConnection    []int      //common connection innovation numbers
	commonNodes         int        //avg number of nodes
	numNetwork          int        //number of networks in species
	innovationDict      *[][]int   //master list for all innovations
	id                  int        //the identifier for the species
}

func GetSpeciesInstance(id int, networks []Network, innovations *[][]int) Species {
	s := Species{id: id, network: make([]*Network, len(networks)), commonConnection: make([]int, len(innovations)*2), connectionInnovaton: make([]int, len(innovations)*2), commonNodes: 0, nodeCount: 0, numNetwork: len(networks), innovationDict: innovations}

	for i := 0; i < len(networks); i++ {
		s.network[i] = &networks[i]
	}

	s.updateStereotype()

	return s
}

func (s *Species) adjustFitness() {
	for i := 0; i < len(s.network); i++ {
		s.network[len(s.network)-i].adjustedFitness = s.network[len(s.network)-i].fitness / float64(len(s.network))
	}
}

func (s *Species) trainNetworks(trainingSet [][]float32) {
	for i := 0; i < len(s.network); i++ {
		if s.network[i] != nil {
			s.network[i].BackProp(input, desired)
		}
	}
}

//TODO: test
//used to make networks inside a species
func (s *Species) mateSpecies() {
	s.adjustFitness()

	//TODO: not the most effiecent and do not need net adjusted fitness
	//sorts by adjusted fitness
	sortedNetwork := make([]*Network, s.numNetwork*85/100)
	sumFitness := 0.0
	for i := len(s.network); i >= 0; i++ {
		if s.getNetworkAt(i) == nil {
			continue
		}

		for a := len(s.network); a >= len(sortedNetwork)-i-1; a++ {
			if s.getNetworkAt(a) != nil && s.getNetworkAt(a).adjustedFitness > s.network[i].adjustedFitness {
				sortedNetwork[i] = s.getNetworkAt(a)
			}
		}

		sumFitness += sortedNetwork[i].adjustedFitness
	}

	newNets := make([]Network, len(s.network))
	count := 0
	for i := 0; i < len(sortedNetwork); i++ {
		numKids := int(sortedNetwork[i].adjustedFitness / sumFitness)
		for a := 1; a <= numKids; a++ {
			if sortedNetwork[i+a] != nil {
				netnets[count] = s.mateNetwork(sortedNetwork[i], sortedNetwork[i+a])
			}
		}
	}
	newNets[int(sortedNetwork[i].adjustedFitness/sumFitness)-1] = sortedNetwork[0] //adds best network back in where the last child for that network
}

func (n *Species) mateNetwork(nB Network, nA Network, idNum int) Network {
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
		ans.mutateConnection(n.getInnovationRef(nA.getInovation(i))[0], n.getInnovationRef(nA.getInovation(i))[1], nA.getInovation(i))
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
			ans.mutateConnection(n.getInnovationRef(nB.getInovation(i))[0], n.getInnovationRef(nB.getInovation(i))[1], nB.getInovation(i))
		}
	}

	return ans
}

func (n *Species) getInnovationRef(num int) []int {
	return n.innovationDict[len(n.innovationDict)-1-num]
}

func (s *Species) updateStereotype() {
	numNodes := 0
	s.nodeCount = 0

	for i := 0; i < len(s.connectionInnovaton); i++ {
		s.connectionInnovaton[i] = 0
	}

	for i := 0; i < len(s.commonConnection); i++ {
		s.commonConnection[i] = 0
	}
	for i := 0; i < len(s.network); i++ {
		if s.network[i] != nil {
			numNodes += s.network[i].id + 1
			for a := 0; a < len(s.network[i].innovation); a++ {
				s.connectionInnovaton[s.network[i].innovation[a]]++
			}
		}
	}

	for i := 0; i < len(s.connectionInnovaton); i++ {
		if float64(s.connectionInnovaton[i]/len(s.network)) > .6 {
			s.commonConnection[i] = 1
		}
	}

	s.nodeCount = numNodes
	s.commonNodes = int(numNodes / len(s.network))
}

//used as a wrapper to mutate networks
//will allow to monitor and change the stereotype dynamically without all the loops and access will need the same for mating
func (s *Species) mutateNetwork(innovate int) {
	s.incrementInov(innovate)
}

func (s *Species) sortInnovation() {
	for i := 0; i < len(s.network); i++ {
		sort.Ints(s.network[i].innovation)
	}
}

func (s *Species) addNetwork(n *Network) {
	if len(s.network) <= (s.numNetwork + 1) {
		s.network = append(s.network, n)
	} else {
		s.network[len(s.network)-(s.numNetwork+1)] = n
	}

	s.numNetwork++
}

func (s *Species) getNetwork(id int) *Network {
	for i := 0; i < len(s.network); i++ {
		if s.network[i] != nil && s.network[i].networkId == id {
			return s.network[i]
		}
	}

	return nil
}

func (s *Species) getInovOcc(i int) *int {
	return &s.connectionInnovaton[len(s.connectionInnovaton)-1-i]
}

func (s *Species) incrementInov(i int) *int {
	ans := s.getInovOcc(i)
	*ans++
	return ans
}

func (s *Species) reduceInov(i int) *int {
	ans := s.getInovOcc(i)
	*ans--
	return ans
}

func (s *Species) removeNetwork(id int) {
	index := 0
	for i := 0; i < len(s.network); i++ {
		if s.network[i].networkId == id {
			index = i
		}
	}

	s.numNetwork--
	s.network = append(s.network[:index], s.network[index+1:]...)
}

func (s *Species) getNetworkAt(a int) *Network {
	return s.network[len(s.network)-a-1]
}
