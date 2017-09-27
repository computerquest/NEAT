package main

import (
	"testing"
)

func TestGetNetworkInstance(t *testing.T) {
	n := GetNetworkInstance(5,5, 0)

	if len(n.nodeList) < 10 {
		t.Errorf("we got", len(n.nodeList))
	}
}
func TestGetNode(t *testing.T) {
	n := GetNetworkInstance(5,5, 0)

	if n.getNode(0).id != 0 {
		t.Errorf("we got", n.getNode(0).id)
	}
}

func TestMutateNode(t *testing.T) {
	n := GetNetworkInstance(5,5, 0)

	n.mutateNode(0, 8)

	ans := false

	for i := 0; i < len(n.nodeList[0].send); i++ {
		if n.nodeList[0].send[i].nodeTo.id == 8{
			ans = true
		}
	}

	if !ans {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", len(n.nodeList), 11)
	}
}
