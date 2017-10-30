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
		n[i] = GetNetworkInstance(5,5, i)
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