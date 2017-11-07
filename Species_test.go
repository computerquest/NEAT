package main

import (
	"testing"
)

/*
tests :
retrieving networks
editing networks
 */
func TestGetSpeciesInstance(t *testing.T) {
	n := make([]Network,10)
	for i := 0; i < len(n); i++ {
		n[i] = GetNetworkInstance(5, 5, i, 0)
	}

	s := GetSpeciesInstance(100, n)

	if s.getNetwork(0).networkId != n[0].networkId || s.getNetwork(1).networkId != n[1].networkId {
		t.Errorf("bad retrieval")
	}

	n[0].fitness = 1000
	if s.getNetwork(0).fitness != n[0].fitness {
		t.Errorf("bad edit")
	}
}

func TestUpdateStereoType(t *testing.T) {
	n := make([]Network, 10)
	for i := 0; i < len(n); i++ {
		n[i] = GetNetworkInstance(5, 5, i, 0)
	}

	s := GetSpeciesInstance(100, n)

	n[0].addInnovation(20)
	n[1].addInnovation(20)
	n[2].addInnovation(20)
	n[3].addInnovation(20)
	n[4].addInnovation(20)
	n[5].addInnovation(20)
	n[6].addInnovation(20)
	n[7].addInnovation(20)

	n[0].addInnovation(30)
	n[1].addInnovation(30)
	n[2].addInnovation(30)
	n[3].addInnovation(30)

	s.updateStereotype()

	if s.commonConnection[20] != 1 || s.commonConnection[30] != 0 {
		t.Errorf("something wrong was included")
	}
}

func TestMutateNetwork(t *testing.T) {
	n := make([]Network, 10)
	for i := 0; i < len(n); i++ {
		n[i] = GetNetworkInstance(5, 5, i, 0)
	}

	s := GetSpeciesInstance(100, n)

	s.mutateNetwork(30)

	if *s.getInovOcc(30) != 1 {
		t.Errorf("didn't increment correctly")
	}
}

func TestGetNetwork(t *testing.T) {
	n := make([]Network, 10)
	for i := 0; i < len(n); i++ {
		n[i] = GetNetworkInstance(5, 5, i, 0)
	}
	s := GetSpeciesInstance(100, n)

	if s.getNetwork(2) != &n[2] {
		t.Errorf("returned the wrong network")
	}
}

func TestRemoveNetwork(t *testing.T) {
	n := make([]Network, 10)
	for i := 0; i < len(n); i++ {
		n[i] = GetNetworkInstance(5, 5, i, 0)
	}
	s := GetSpeciesInstance(100, n)

	s.removeNetwork(2)

	if s.getNetwork(2) != nil {
		t.Errorf("didn't remove correctly")
	}
}