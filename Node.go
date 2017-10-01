package main

import (
	"math"
)

type Node struct {
	value float64
	id int
	receive []*Connection //connections to this node
	send []Connection //connections sent from this node
	influence float64 //this nodes influence (used for backprop)
	inputRecieved int //number of connections that have responded with input values
	influenceRecieved int //number of connections that have responded with influence values
	//activation bool //used to signal input nodes don't need activation but might not need
	numConIn int
	numConOut int
}

//calculate input to this node
func (n Node) netInput() float64 {
	var sum float64 = 0
	for i := 0; i < len(n.receive); i++ {
		c := n.receive[i]
		if !c.disable {
			sum += (c.nodeFrom.value)*c.weight
		}
	}

	return sum
}

//called when connection recieves a input value
func (n *Node) recieveValue() {
	n.inputRecieved++

	if n.inputRecieved == len(n.receive) {
		n.setValue(sigmoid(n.netInput()))
		n.inputRecieved = 0
	}
}

//called when connection recieves an influence value
func (n *Node) recieveInfluence() {
	n.influenceRecieved++

	if n.influenceRecieved == len(n.send) {
		n.influence = 0
		for i := 0; i < len(n.send); i++ {
			sumInfluence := 0.0
			if ! n.send[i].disable {
				sumInfluence += n.send[i].nodeTo.influence * n.send[i].weight
			}

			n.setInfluence(sumInfluence)
		}

		n.influenceRecieved = 0
	}
}

func (n *Node) setValue(i float64) {
	n.value = i
	n.signalValue()
}

func (n *Node) setInfluence(i float64) {
	n.influence = i
	n.signalInfluence()
}

//notifies all connections that the value has been calculated
func (n *Node) signalValue() {
	for i := 0; i < len(n.send); i++ {
		n.send[i].notifyValue()
	}
}

//notifies all connections that the influence has been calculated
func (n *Node) signalInfluence() {
	for i := 0; i < len(n.receive); i++ {
		n.receive[i].notifyInfluence()
	}
}

func sigmoid(value float64) float64 {
	return 1 / (1 + (1/math.Pow(2.71, value)))
}

func sigmoidDerivative(value float64) float64 {
	return sigmoid(value)*(1 - sigmoid(value))
}

/*
could have it so that the add methods will add the pointer to themselves
 */
func (n *Node) addSendCon(c Connection) *Connection {
	if len(n.send)-1 > n.numConOut {
		n.send[len(n.send)-n.numConOut-1] = c
		n.numConOut++
	} else {
		n.send = append(n.send, c)
	}

	return &n.send[len(n.send)-n.numConOut-2]
}
func (n *Node) addRecCon(c *Connection) *Connection{
	if len(n.receive)-1 > n.numConIn {
		n.receive[len(n.receive)-n.numConIn-1] = c
		n.numConIn++
	} else {
		n.receive = append(n.receive, c)
	}

	return c
}
func (n *Node) getRecCon(i int) *Connection {
	return n.receive[len(n.receive)-1-i]
}

func (n *Node) getSendCon(i int) *Connection {
	return &n.send[len(n.send)-1-i]
}