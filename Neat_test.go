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
