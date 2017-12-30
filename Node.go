package main

import (
	"math"
)

//MAX 100 CONNECTIONS

type Node struct {
	value             float64
	id                int
	receive           []*Connection //pointer to connections sent to this node
	send              []Connection  //connections sent from this node
	influence         float64       //this nodes influence (used for backprop)
	inputRecieved     int           //number of connections that have responded with input values
	influenceRecieved int           //number of connections that have responded with influence values
}

/////////////////////////////////////////////PASS
//recieves value from recieving connections and can set value when ready
func (n *Node) recieveValue() {
	n.inputRecieved++

	if n.inputRecieved == len(n.receive) {
		sum := 0.0
		for i := 0; i < len(n.receive); i++ {
			c := n.receive[i]
			if !c.disable {
				sum += (c.nodeFrom.value) * c.weight
			}
		}

		n.setValue(tanh(sum))
		n.inputRecieved = 0
	}
}

//recieves influence from sending connections and can set influence when ready
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

/////////////////////////////////////////////CONNECTION
func (n *Node) addSendCon(c Connection) *Connection {
	if len(n.send) >= cap(n.send) {
		n.send = append(n.send, c)
	} else {
		n.send = n.send[0: len(n.send)+1]
		n.send[len(n.send)-1] = c
	}

	return &n.send[len(n.send)-1]
}
func (n *Node) addRecCon(c *Connection) *Connection {
	if len(n.receive) >= cap(n.receive) {
		n.receive = append(n.receive, c)
	} else {
		n.receive = n.receive[0: len(n.receive)+1]
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
//sets value and signals to higher (send connection) nodes value is calculated
func (n *Node) setValue(i float64) {
	n.value = i

	for i := 0; i < len(n.send); i++ {
		n.send[i].notifyValue()
	}
}

//sets influence and signals to lower (recieve connection) nodes influence is calculated
func (n *Node) setInfluence(i float64) {
	n.influence = i * tanhDerivative(n.value)
	for i := 0; i < len(n.receive); i++ {
		if n.receive[i] != nil {
			n.receive[i].notifyInfluence()
		}
	}
}

//////////////////////////////////////////////ACTIVATION
func tanh(value float64) float64 {
	return (math.Pow(2.71, value) - math.Pow(2.71, -1*value)) / (math.Pow(2.71, value) + math.Pow(2.71, -1*value))
}
func tanhDerivative(value float64) float64 {
	return 1 - math.Pow(tanh(value), 2)
}
func sigmoid(value float64) float64 {
	return 1 / (1 + (1 / math.Pow(2.71, value)))
}
func sigmoidDerivative(value float64) float64 {
	return sigmoid(value) * (1 - sigmoid(value))
}

///////////////////////////////////////TYPE
//determines if a node is input (bias nodes will evaluate true)
func isInput(n *Node) bool {
	if cap(n.receive) == 0 {
		return true
	}

	return false
}

//determines if output node
func isOutput(n *Node) bool {
	if cap(n.send) == 0 {
		return true
	}

	return false
}
