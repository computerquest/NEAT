package main

import (
	"sort"
)

type Species struct {
	network             []*Network
	connectionInnovaton []int //this will be the max size of innovation number
	nodeCount           int
	commonConnection    []int
	commonNodes         int
	numNetwork 			int
}

func GetSpeciesInstance(maxInnovation int, networks []Network) Species {
	s := Species{network: make([]*Network, cap(networks)), commonConnection: make([]int, int(maxInnovation*2)), connectionInnovaton: make([]int, int(maxInnovation*2)), commonNodes: 0, nodeCount: 0}

	for i := 0; i < len(networks); i++ {
		s.network[i] = &networks[i]
	}

	//todo uncomment after testing
	//s.updateStereotype()

	return s
}

func (s *Species) adjustFitness() {
	for i := 0; i < len(s.network); i++ {
		s.network[len(s.network)-i].adjustedFitness = s.network[len(s.network)-i].fitness / float64(len(s.network))
	}
}

//todo have an add max innovation method
//todo make sure it adds innovations upon creation
func (s *Species) mate(n *Network, nA *Network) Network{
	s.numNetwork++
	newNetwork := *n
	newNetwork.networkId = s.numNetwork

	for nA.id > newNetwork.id {
		newNetwork.createNode()
	}

	//todo simplify
	for i := 0; i <= nA.id; i ++ {
		node := nA.getNode(i)
		for a := 0; a < len(node.send); a++ {
			checkNum := node.getSendCon(a).inNumber
			contains := false
			for b := 0; b < len(newNetwork.innovation); b++ {
				if newNetwork.innovation[b] == checkNum {
					contains = true
					break
				}
			}

			if !contains {
			}
		}
	}

	return newNetwork
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
		numNodes += s.network[i].id+1
		for a := 0; a < len(s.network[i].innovation); a++ {
			s.connectionInnovaton[s.network[i].innovation[a]]++
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
	if len(s.network) <= (s.numNetwork+1) {
		s.network = append(s.network, n)
	} else {
		s.network[len(s.network)-(s.numNetwork+1)] = n
	}

	s.numNetwork++
}

func (s *Species) getNetwork(id int) *Network {
	for i := 0; i < len(s.network); i++ {
		if s.network[i].networkId == id {
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

	s.network = append(s.network[:index], s.network[index+1:]...)
}