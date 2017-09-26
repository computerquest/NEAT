package main

import (
	"testing"
)

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
