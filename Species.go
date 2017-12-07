package main

import (
	"sort"
)

//TODO: might want to consider starting the innovation master list at one so that all of the arrays have a default value (or prevents default value)
//TODO: look into node id system and make sure that it doesn't allow different types of nodes to have the same id (screw innovation number pairings)
//TODO: fix the avg because empty slots created by append will screw
//TODO: make sure when mate change the neat class networks
//TODO: make sure that length is always exact
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
	s := Species{id: id, network: make([]*Network, len(networks)), commonConnection: make([]int, len(*innovations)*2), connectionInnovaton: make([]int, len(*innovations)*2), commonNodes: 0, nodeCount: 0, numNetwork: len(networks), innovationDict: innovations}

	for i := 0; i < len(networks); i++ {
		s.network[i] = &networks[i]
	}

	s.updateStereotype()

	return s
}

func (s *Species) adjustFitness() {
	for i := 0; i < len(s.network); i++ {
		s.network[i].adjustedFitness = s.network[i].fitness / float64(len(s.network))
	}
}

func (s *Species) trainNetworks(trainingSet [][][]float64) {
	for i := 0; i < len(s.network); i++ {
		if s.network[i] != nil {
			s.network[i].trainSet(trainingSet, 1500)
		}
	}
}

//used to make networks inside a species
func (s *Species) mateSpecies() []Network {
	s.adjustFitness()

	//TODO: not the most effiecent and do not need net adjusted fitness
	//sorts by adjusted fitness
	sortedNetwork := make([]*Network, s.numNetwork*85/100)
	lastValue := 0.0
	sumFitness := 0.0
	for i := 0; i < len(sortedNetwork); i++ { //TODO: why
		if s.getNetworkAt(i) == nil {
			continue
		}

		localMax := 0.0
		localIndex := 0
		for a := 0; a < len(s.network); a++ {
			if s.getNetworkAt(a) != nil && s.getNetworkAt(a).adjustedFitness > localMax && s.getNetworkAt(a).adjustedFitness < lastValue {
				localMax = s.network[a].adjustedFitness
				localIndex = a
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
		for a := 1; a <= numKids && a+i < len(sortedNetwork); a++ {
			if sortedNetwork[i+a] != nil {
				newNets[count] = s.mateNetwork(*sortedNetwork[i], *sortedNetwork[i+a])
				count++
				numMade--
			}
		}

		//TODO: finish so that all the kids are made
		/*for numMade > 0 {
			newNets[count] = s.mateNetwork(*sortedNetwork[i], *sortedNetwork[i+a])
		}*/
	}
	newNets[int(sortedNetwork[0].adjustedFitness/sumFitness*float64(len(newNets))-float64(1))] = *sortedNetwork[0] //adds best network back in where the last child for that network

	for i := 0; i < len(newNets); i++ {
		newNets[i].networkId = s.network[i].networkId
	}

	s.updateStereotype()

	return newNets
}

func isRealSpecies(s *Species) bool {
	if cap(s.network) != 0 {
		return true
	}
	return false
}
func (n *Species) mateNetwork(nB Network, nA Network) Network {
	ans := GetNetworkInstance(len(nB.output), len(nB.input), 0, nB.species, .1)

	var numNode int
	if nA.id > nB.id {
		numNode = nA.id
	} else {
		numNode = nB.id
	}

	for i := ans.id; i < numNode; i++ {
		ans.createNode()
	}

	for i := 0; i < len(nA.innovation); i++ {
		ans.mutateConnection(n.getInnovationRef(nA.getInovation(i))[0], n.getInnovationRef(nA.getInovation(i))[1], nA.getInovation(i))
	}

	for i := 0; i < len(nB.innovation); i++ {
		exist := false
		for a := 0; a < len(nA.innovation); a++ {
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
	return (*n.innovationDict)[num]
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
				s.incrementInov(s.network[i].innovation[a])
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
	if len(s.network) >= cap(s.network) {
		s.network = append(s.network, n)
	} else {
		s.network = s.network[0 : len(s.network)+1]
		s.network[len(s.network)-1] = n
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
	if i >= len(s.connectionInnovaton) {
		insert := make([]int, 1+i-len(s.connectionInnovaton))
		s.connectionInnovaton = append(s.connectionInnovaton, insert...)
		s.commonConnection = append(s.commonConnection, insert...)
	}
	return &s.connectionInnovaton[i]
}

func (s *Species) incrementInov(i int) *int {
	ans := s.getInovOcc(i)
	*ans++

	if float64(s.connectionInnovaton[i]/len(s.network)) > .6 {
		s.commonConnection[i] = 1
	}
	return ans
}

func (s *Species) reduceInov(i int) *int {
	ans := s.getInovOcc(i)
	*ans--

	if float64(s.connectionInnovaton[i]/len(s.network)) <= .6 {
		s.commonConnection[i] = 0
	}

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
	return s.network[a]
}
