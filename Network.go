package main

import (
	"fmt"
)

//NOTE most of the calculating work is networked by nodes inside the struct

type Network struct {
	nodeList []Node //master list of nodes
	numConnections int
	numNodes int
	innovation []int //list of inovation numbers this network has
	id int //network id
	learningRate float64 //learning rate for backprop
	output []*Node //output nodes
	input []*Node //input nodes
}

//processes the network
func (n *Network) Process(input []float64) {
	for i := 0; i < len(n.input); i++ {
		n.input[i].setValue(input[i])
	}
}

//todo test when time
//backpropogates the network to desired one time
func (n *Network) BackProp(input []float64, desired []float64) {
	n.Process(input) //need to do so that you are performing the algorithm on that set of values

	//this will calc all the influence
	for i := 0; i < len(n.output); i++ {
		n.output[i].setInfluence(n.output[i].value-desired[i])
	}

	//actually adjusts the weights
	for i := 0; i < len(n.nodeList); i++ {
		derivative := sigmoidDerivative(n.nodeList[i].value)
		for a := 0; a < len(n.nodeList[i].receive); a++ {
			n.nodeList[i].receive[a].nextWeight +=  derivative * (n.nodeList[i].receive[a].nodeFrom.value) * n.nodeList[i].influence * n.learningRate
		}
	}
	//backprop output and hidden
	/*for z := 2; z >= 1; z++ {
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
	}*/
}

//todo test
func (n *Network) addConnection(from int, to int) {
	n.nodeList[from].send = append(n.nodeList[from].send,   Connection{weight: 1, nextWeight: 0, disable:false, nodeFrom: &n.nodeList[from], nodeTo: &n.nodeList[to]})
	n.nodeList[to].receive = append(n.nodeList[to].receive, &n.nodeList[from].send[len( n.nodeList[from].send)-1])
}

//todo finish
func (n *Network) addNode(from int, to int) {

}

//todo need to make sure doing the right connections
func (n *Network) GetInstance(input int, output int) {
	//set all default values
	n.learningRate = .1
	count := 0
	n.numConnections = 0
	n.numNodes = input+output

	n.nodeList = make([]Node, (input+output)*2)
	n.output = make([]*Node, output)
	n.input = make([]*Node, input)

	fmt.Print("initialized")

	//create output nodes
	for i := 0; i < output; i++ {
		n.nodeList[count] = Node {value:0, influenceRecieved: 0, inputRecieved: 0, id:n.id, receive:make([]*Connection, input)}
		n.output[i] = &n.nodeList[count]
		count++
		n.id++
	}
	fmt.Print("output")

	//creates the input nodes and adds them to the network
	for i := 0; i < input; i++ {
		n.nodeList[count] = Node {value:0, id:n.id, influenceRecieved: 0, inputRecieved: 0, send:make([]Connection, output)}
		n.input[i] = &n.nodeList[count]

		//creates the connections
		for a := 0; a < output; a++ {
			n.nodeList[count].send[a] = Connection{weight: 1, nextWeight: 0, disable:false, nodeFrom: n.input[i], nodeTo: n.output[a]}
			n.nodeList[a].receive[i] = &n.nodeList[count].send[a]
			n.numConnections++
		}

		n.id++
		count++
	}
	fmt.Print("input")
}
