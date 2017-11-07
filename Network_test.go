package main

import (
	"testing"
)

func TestGetNetworkInstance(t *testing.T) {
	n := GetNetworkInstance(5, 5, 0, 0)

	if n.numConnections < 10 || n.id < 9 {
		t.Errorf("we got %d and %d", n.numConnections, n.id)
	}
}

func TestCreateNode(t *testing.T) {
	n := GetNetworkInstance(5, 5, 0, 0)

	a := n.createNode()
	b := n.createNode()
	if a.id == b.id  || a.id != 10 || b.id != 11 {
		t.Errorf("we got %d and %d", a.id, b.id)
	}
}
func TestGetNode(t *testing.T) {
	n := GetNetworkInstance(5, 5, 0, 0)

	if n.getNode(0).id != 0 || n.getNode(9).id != 9 {
		t.Errorf("we got", n.getNode(0).id)
	}
}

func TestMutateNode(t *testing.T) {
	n := GetNetworkInstance(5, 5, 0, 0)

	num := n.mutateNode(5, 0, 10, 11)

	ans := true
	ansA := false
	ansB := false
	//the node to is nil (because this has default initial
	for i := 0; i < len(n.getNode(0).send); i++ {
		if n.getNode(0).getSendCon(i).nodeTo != nil && n.getNode(0).getSendCon(i).nodeTo.id == 0 {
			ans = false
		} else if n.getNode(0).getSendCon(i).nodeTo != nil && n.getNode(0).getSendCon(i).nodeTo.id == 0 {
			ansA = true

			for a := 0; a < len(n.getNode(0).getSendCon(i).nodeTo.send); a++ {
				if n.getNode(0).getSendCon(i).nodeTo.getSendCon(a).nodeTo.id == 5 {
					ansB = true
				}
			}
		}
	}

	if !ans && ansA && ansB && n.getNode(num).id ==10 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", ans, ansA, n.getNode(num).id)
	}
}

func TestMutateConnection(t *testing.T) {
	n := GetNetworkInstance(5, 5, 0, 0)
	num := n.mutateNode(5,0, 100, 101)

	n.mutateConnection(num, 9, 1000)

	t.Log(n.getNode(num))
	if n.getNode(num).id != num  || n.getNode(5).id != 5{
		t.Error("we got the wrong node")
	}

	ans := false
	for i := 0; i < len(n.getNode(num).send); i++ {
		t.Log(n.getNode(num).send[i])

		if nil != n.getNode(num).send[i].nodeTo{
			t.Log(n.getNode(num).send[i].nodeTo.id)
		}

		if nil != n.getNode(num).send[i].nodeTo && n.getNode(num).send[i].nodeTo.id == 0 {
			ans = true
		}
	}

	if !ans {
		t.Errorf("Not found")
	}
}

func TestAddInnovation(t *testing.T) {
	n := GetNetworkInstance(5, 5, 0, 0)

	if n.innovation[1] != 1 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", 1, n.innovation[1])
	}
}

func TestRemoveInnovation(t *testing.T) {
	n := GetNetworkInstance(5, 5, 0, 0)

	n.removeInnovation(1)

	if n.innovation[1] == 1 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", 1, n.innovation[1])
	}
}
