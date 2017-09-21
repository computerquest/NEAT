package main

import "sort"

type Species struct {
	network             []*Network
	connectionInnovaton []int //this will be the max size of innovation number
	nodeCount           int
	commonConnection    []int
	commonNodes         int
}

func GetSpeciesInstance(maxInnovation int, networks []*Network) Species {
	s := Species{network: make([]*Network, cap(networks)), connectionInnovaton: make([]int, int(maxInnovation*1.25)), commonNodes: 0, nodeCount: 0,}

	//doing this so slice passed is not kept in memory
	for i := 0; i < len(s.network); i++ {
		s.network[i] = networks[i]
	}

	s.updateStereotype()

	return s
}

func (s *Species) adjustFitness() {
	for i := 0; i < len(s.network); i++ {
		s.network[i].adjustedFitness = s.network[i].fitness / float64(len(s.network))
	}
}

//todo finish
func (s *Species) mate() {

}

//todo test
func (s *Species) updateStereotype() {
	numNodes := 0
	s.nodeCount = 0

	for i := 0; i < len(s.connectionInnovaton); i++ {
		s.connectionInnovaton[i] = 0
	}

	for i := 0; i < len(s.network); i++ {
		numNodes += s.network[i].numNodes
		for a := 0; a < len(s.network[i].innovation); a++ {
			if s.network[i].innovation[a] >= len(s.connectionInnovaton) {
				s.connectionInnovaton = append(s.connectionInnovaton)
			}
			s.connectionInnovaton[s.network[i].innovation[a]]++
		}
	}

	count := 0
	for i := 0; i < len(s.connectionInnovaton); i++ {
		if s.connectionInnovaton[i]/len(s.network) > .6 {
			s.commonConnection[count] = s.connectionInnovaton[i]
		}
	}

	s.nodeCount = numNodes
	s.commonNodes = int(numNodes / len(s.network))
}

//todo finish
//used as a wrapper to mutate networks
//will allow to monitor and change the stereotype dynamically without all the loops and access will need the same for mating
func (s *Species) mutateNetwork() {

}

func (s *Species) sortInnovation() {
	for i := 0; i < len(s.network); i++ {
		sort.Ints(s.network[i].innovation)
	}
}
