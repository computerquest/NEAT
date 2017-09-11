package main

import "math"

type Node struct {
	value float64
	id int
	receive []*Connection
	send []Connection //this list is seqential for initialization
	influence float64
	inputRecieved int
	influenceRecieved int
	//activation bool //used to signal input nodes don't need activation but might not need
}

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

func (n *Node) signalValue() {
	for i := 0; i < len(n.send); i++ {
		n.send[i].notifyValue()
	}
}

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
