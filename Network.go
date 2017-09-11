package main

import (
	"fmt"
	"math"
)

type Network struct {
	//nodes [][]*Node // this is used for the backprop and processing
	nodeList []Node //this is used primarily for connections and mating
	numConnections int
	numNodes int
	innovation []int
	id int
	learningRate float64
	output []*Node
	input []*Node
}

//todo finish
/*
could make it so that connections hold pointers to nodes
right now it looks that everynode will call every other node up
we would do the same but backwards with backprop
to solve the waiting issue I would have the nodes wait until the same amount of sending or recieving nodes have been activated
 */
func (n *Network) Process(input []float64) {
	for i := 0; i < len(n.input); i++ {
		n.input[i].value = input[i];
	}

	for i := 1; i < len(n.nodes); i++ {
		for a := 0; a < len(n.nodes[i]); a++ {
			n.nodes[i][a].value = sigmoid(n.nodes[i][a].netInput())
		}
	}
}

//todo finish/test
func (n *Network) BackProp(input []float64, desired []float64) {
	n.Process(input)

	//backprop output and hidden
	for z := 2; z >= 1; z++ {
		for i := 0; i < len(n.nodes[z]); i++ {
			node := n.nodes[z][i]

			node.influence = 0
			derivative := sigmoidDerivative(node.value)

			if z < 2 {
				for a := 0; a < len(node.receive); a++ {
					node.influence += (*node.receive[a].connectInfluence) * (node.receive[a].weight)
				}
			} else {
				node.influence = node.value-desired[i]
			}

			for a := 0; a < len(node.receive); a++ {
				node.receive[a].nextWeight += derivative * (*node.receive[a].sendValue) * node.influence * n.learningRate
			}
		}
	}
}

//todo finish
func (n *Network) addConnection(from int, to int) {
	n.nodeList[from].send = append(n.nodeList[from].send,  Connection{weight: 1, disable:false, sendValue:&n.nodeList[from].value, connectInfluence:&n.nodeList[to].influence})
	n.nodeList[to].receive = append(n.nodeList[to].receive, &n.nodeList[from].send[len( n.nodeList[from].send)-1])

}

//todo finish
func (n *Network) addNode(from int, to int) {

}

func sigmoid(value float64) float64 {
	return 1 / (1 + (1/math.Pow(2.71, value)))
}

func sigmoidDerivative(value float64) float64 {
	return sigmoid(value)*(1 - sigmoid(value))
}

//todo need to make sure doing the right connections and test
func (n *Network) GetInstance(input int, output int) {
	n.learningRate = .1
	count := 0
	n.numConnections = 0
	n.numNodes = 0

	n.nodeList = make([]Node, (input+output)*2)

	fmt.Print("initialized")

	for i := 0; i < output; i++ {
		n.nodeList[count] = Node {value:0, id:n.id, receive:make([]*Connection, input)}
		n.output[i] = &n.nodeList[count]
		count++
		n.id++
	}
	fmt.Print("output")

	//creates the nodes and adds them to the network
	for i := 0; i < input; i++ {
		n.nodeList[count] = Node {value:0, id:n.id, send:make([]Connection, output)}
		n.output[i] = &n.nodeList[count]

		//creates the connections
		for a := 0; a < output; a++ {
			n.nodeList[count].send[a] = Connection{weight: 1, disable:false, idFrom: n.nodeList[count].id, idTo: n.nodeList[a].id}
			n.nodeList[a].receive[i] = &n.nodeList[count].send[a]
			n.numConnections++
		}

		n.id++
		count++
	}
	fmt.Print("input")

	n.numNodes = input+output
}
