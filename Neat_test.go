package main

import "testing"

func TestGetNeatInstance(t *testing.T) {
	s := GetNeatInstance(40, 5, 5)

	if &s.species[0].network[len(s.species[0].network)-1] == &s.species[1].network[0] {
		t.Errorf("networks double included")
	}
	if len(s.network) != 40 {
		t.Errorf("wrong number of networks")
	}
}

func TestMutatePopulation(t *testing.T) {
	s := GetNeatInstance(40, 5, 5)

	og := s.species[0].getNetworkAt(0).numInnovation

	s.mutatePopulationTest()

	if s.species[0].getNetworkAt(0).numInnovation == og {
		t.Errorf("didnt do anything")
	}
}

func TestMateNetwork(t *testing.T) {
	s := GetNeatInstance(40, 5, 5)

	if len(s.network) != 40 {
		t.Errorf("didnt do anything")
	}
	if len(s.network) != 40 {
		t.Errorf("bad new innovation")
	}
}
