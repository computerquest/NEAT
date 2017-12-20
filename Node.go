package main

import (
	"math"
)

//MAX 100 CONNECTIONS

//TODO: can't give bias node connections because it is seen as an input and output in the classification methods
type Node struct {
	value             float64
	id                int
	receive           []*Connection //connections to this node
	send              []Connection  //connections sent from this node
	influence         float64       //this nodes influence (used for backprop)
	inputRecieved     int           //number of connections that have responded with input values
	influenceRecieved int           //number of connections that have responded with influence values
}

/////////////////////////////////////////////PASS
func (n *Node) recieveValue() {
	n.inputRecieved++

	if n.inputRecieved == len(n.receive) {
		n.setValue(sigmoid(n.netInput()))
		n.inputRecieved = 0
	}
}
func (n *Node) recieveInfluence() {
	n.influenceRecieved++

	if n.influenceRecieved == len(n.send) {
		n.influence = 0
		for i := 0; i < len(n.send); i++ {
			if !n.send[i].disable {
				n.influence += n.send[i].nodeTo.influence * n.send[i].weight
			}

		}
		n.setInfluence(n.influence)
		n.influenceRecieved = 0
	}
}
func (n *Node) signalValue() {
	for i := 0; i < len(n.send); i++ {
		n.send[i].notifyValue()
	}
}
func (n *Node) signalInfluence() {
	for i := 0; i < len(n.receive); i++ {
		if n.receive[i] != nil {
			n.receive[i].notifyInfluence()
		}
	}
}

/////////////////////////////////////////////CONNECTION
func (n *Node) addSendCon(c Connection) *Connection {
	if len(n.send) >= cap(n.send) {
		n.send = append(n.send, c)
	} else {
		n.send = n.send[0 : len(n.send)+1]
		n.send[len(n.send)-1] = c
	}

	return &n.send[len(n.send)-1]
}
func (n *Node) addRecCon(c *Connection) *Connection {
	if len(n.receive) >= cap(n.receive) {
		n.receive = append(n.receive, c)
	} else {
		n.receive = n.receive[0 : len(n.receive)+1]
		n.receive[len(n.receive)-1] = c
	}

	return c
}
func (n *Node) getRecCon(i int) *Connection {
	return n.receive[i]
}
func (n *Node) getSendCon(i int) *Connection {
	return &n.send[i]
}

/////////////////////////////////////////////////PROCESS
func (n Node) netInput() float64 {
	var sum float64 = 0
	for i := 0; i < len(n.receive); i++ {
		c := n.receive[i]
		if !c.disable {
			sum += (c.nodeFrom.value) * c.weight
		}
	}

	return sum
}
func (n *Node) setValue(i float64) {
	n.value = i
	n.signalValue()
}
func (n *Node) setInfluence(i float64) {
	n.influence = i
	n.signalInfluence()
}
func sigmoid(value float64) float64 {
	return 1 / (1 + (1 / math.Pow(2.71, value)))
}
func sigmoidDerivative(value float64) float64 {
	return sigmoid(value) * (1 - sigmoid(value))
}

//TODO: do i need any of these node type methods
func isInput(n *Node) bool {
	if cap(n.receive) == 0 {
		return true
	}

	return false
}
func isOutput(n *Node) bool {
	if cap(n.send) == 0 {
		return true
	}

	return false
}

func isBias(n *Node) bool {
	if len(n.send) == 0 {
		return true
	}

	return false
}
