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
	connectionInnovaton []int      //holds number of occurrence of each innovation
	commonInnovation    []int      //common connection innovation numbers
	innovationDict      *[][]int   //master list for all innovations (innovation: 0, nodeFrom, 2, nodeTo, 0)
	id                  int        //the identifier for the species
	mutate              float64    //the odds for a node mutate (ie: .3)
}

var dictControl sync.Mutex //controls access to modifying the connection dictionary
//var dictControlB sync.Mutex //controls access to reading the connection dictionary
//var overall sync.Mutex      //controls access to reading the connection dictionary

func GetSpeciesInstance(id int, networks []Network, innovations *[][]int, mutate float64) Species {
	s := Species{mutate: mutate, id: id, network: make([]*Network, len(networks)), commonInnovation: make([]int, 0, len(*innovations)*2), connectionInnovaton: make([]int, len(*innovations)*2), innovationDict: innovations}

	for i := 0; i < len(networks); i++ {
		s.network[i] = &networks[i]
	}

	s.updateStereotype()

	return s
}

////////////////////////////////////////////////////////////INNOVATION
//CI = commonInnovation
func (s *Species) addCI(a int) {
	for i := 0; i < len(s.commonInnovation); i++ {
		if s.commonInnovation[i] == a {
			return
		}
	}

	if len(s.commonInnovation) >= cap(s.commonInnovation) {
		s.commonInnovation = append(s.commonInnovation, a)
	} else {
		s.commonInnovation = s.commonInnovation[0 : len(s.commonInnovation)+1]
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

	if float64(*ans)/float64(len(s.network)) >= .5 {
		s.addCI(i)
	}

	return ans
}
func (s *Species) reduceInov(i int) *int {
	ans := s.getInovOcc(i)
	*ans--

	if float64(*ans)/float64(len(s.network)) < .5 {
		s.removeCI(i)
	}

	return ans
}

//rechecks all connectionInnovation for commonInnovation
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

//get innovation from innovationDict
func (n *Species) getInnovationRef(num int) []int {
	//overall.Lock()
	//dictControlB.Lock()
	//defer dictControlB.Unlock()
	//defer overall.Unlock()
	//dictControl.Lock()
	//defer dictControl.Unlock()
	return (*n.innovationDict)[num]
}

//adds all new innovation to innovationDict
func (n *Species) createNewInnovation(values []int) int {
	//overall.Lock()
	dictControl.Lock()
	if len(*n.innovationDict) == cap(*n.innovationDict) {
		(*n.innovationDict) = append(*n.innovationDict, values)
	} else {
		*n.innovationDict = (*n.innovationDict)[0 : len(*n.innovationDict)+1]
		(*n.innovationDict)[len(*n.innovationDict)-1] = values
	}
	defer dictControl.Unlock()
	//defer overall.Unlock()

	return len(*n.innovationDict) - 1
}

//sorts all the networks innovation
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
		s.network = s.network[0 : len(s.network)+1]
		s.network[len(s.network)-1] = n
	}

	n.species = s.id

	for i := 0; i < len(n.innovation); i++ {
		s.incrementInov(n.innovation[i])
	}

	s.checkCI()
}

///////////////////////////////////////////////////////////MATE+MUTATE
//recounts all connectionInnovation
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
func (s *Species) mutateNetwork(network *Network) {
	nodeRange := len(network.nodeList)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	//finds or adds innovation numbers and returns the innovation
	addConnectionInnovation := func(numFrom int, numTo int) int {
		//checks to see if preexisting innovation
		maxPos := len(*s.innovationDict)
		for i := 0; i < maxPos; i++ {
			pos := s.getInnovationRef(i)
			if pos[1] == numTo && pos[0] == numFrom {
				s.incrementInov(i)

				return i
			}
		}

		//checks to see if needs to grow
		num := s.createNewInnovation([]int{numFrom, numTo})

		s.incrementInov(num)

		return num
	}

	//mutates a node
	nodeMutate := func() {
		var firstNode int
		var secondNode int
		ans := false

		//picks a node
		for !ans {
			firstNode = int(rand.Float64() * float64(nodeRange))

			if len(network.getNode(firstNode).send) > 0 {
				ans = true
			}
		}

		//picks a random connection from firstNode and gets the id
		secondNode = network.getNode(firstNode).send[int(rand.Float64()*float64(len(network.getNode(firstNode).send)))].nodeTo.id //int(r.Int63n(int64(nodeRange)))

		network.mutateNode(firstNode, secondNode, addConnectionInnovation(firstNode, network.getNextNodeId()), addConnectionInnovation(network.getNextNodeId(), secondNode))
	}

	//randomly picks if node or connection mutate
	if r.Float64() <= s.mutate {
		nodeMutate()
	} else {
		var firstNode int
		var secondNode int
		//ans := true
		attempts := 0 //attempts at finding nodes

		//find 2 unconnected nodes
		//for ans && attempts <= 10 {
		for attempts <= 10 {
			firstNode = int(r.Int63n(int64(nodeRange)))
			secondNode = int(r.Int63n(int64(nodeRange)))

			if firstNode == secondNode || isOutput(network.getNode(firstNode)) || isInput(network.getNode(secondNode)) {
				continue
			}

			if network.getNode(firstNode).connectsTo(secondNode) || network.getNode(secondNode).connectsTo(firstNode) || network.checkCircleMaster(network.getNode(firstNode), secondNode) {
				attempts++
				continue
			}
		}

		//muateNode if attempts exceeded else add the connection
		if attempts > 10 {
			nodeMutate()
		} else {
			network.mutateConnection(firstNode, secondNode, addConnectionInnovation(firstNode, secondNode))
		}
	}
}
func (s *Species) mateNetwork(nB Network, nA Network) Network {
	ans := GetNetworkInstance(len(nB.input)-1, len(nB.output), 0, nB.species, nB.learningRate, false) //child

	numNode := -1 * (len(nB.output) + len(nB.input)) //subtract input and output nodes because those are already created
	if len(nA.nodeList) > len(nB.nodeList) {
		numNode += len(nA.nodeList)
	} else {
		numNode += len(nB.nodeList)
	}

	//create the nodes
	for i := 0; i < numNode; i++ { //this should be ok
		ans.createNode(100)
	}

	//add nA innovation
	for i := 0; i < len(nA.innovation); i++ {
		in := s.getInnovationRef(nA.getInovation(i))
		ans.mutateConnection(in[0], in[1], nA.getInovation(i))
	}

	//add unique nB innovation
	for i := 0; i < len(nB.innovation); i++ {
		firstNode := (*s.innovationDict)[nB.innovation[i]][0]
		secondNode := (*s.innovationDict)[nB.innovation[i]][1]

		//checks to make sure their is no conflict in possible innovations
		if !ans.containsInnovation(nB.innovation[i]) && !(ans.getNode(firstNode).connectsTo(secondNode) || ans.getNode(secondNode).connectsTo(firstNode)) && !ans.checkCircleMaster(ans.getNode(firstNode), secondNode) {
			in := s.getInnovationRef(nB.getInovation(i))
			ans.mutateConnection(in[0], in[1], nB.getInovation(i))
		}
	}

	return ans
}

//trains networks and signals the waitgroup when done
func (s *Species) trainNetworks(trainingSet [][][]float64, control *sync.WaitGroup) {
	for i := 0; i < len(s.network); i++ {
		s.network[i].trainSet(trainingSet, 10000) //I have capped the number of interations intentionaly to control training time
	}
	control.Done()
}

//used to make networks inside a species returns the new networks
func (s *Species) mateSpecies(wg *sync.WaitGroup) {
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
	//mates networks
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

	//mutates for remainder of spots available
	for i := 0; count < len(newNets); i++ {
		newNets[count] = clone(sortedNetwork[i])
		s.mutateNetwork(&newNets[count])
		count++

		if i == len(sortedNetwork)-1 {
			i-- //this can lead to mutating the same network as last time (stacking mutations) but i don't think it is a big deal
		}
	}

	//sets the id
	for i := 0; i < len(newNets); i++ {
		newNets[i].networkId = s.network[i].networkId
		*s.network[i] = newNets[i]
	}

	s.updateStereotype()

	wg.Done()
}

//adjust networks fitness acording to the sum of the species
func (s *Species) adjustFitness() {
	for i := 0; i < len(s.network); i++ {
		s.network[i].adjustedFitness = s.network[i].fitness / float64(len(s.network))
	}
}

//returns the average number of nodes a network has
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
