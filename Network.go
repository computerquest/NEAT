package main

import (
	"fmt"
)

//NOTE most of the calculating work is networked by nodes inside the struct

type Network struct {
	nodeList []Node //master list of nodes
	numConnections int
	innovation []int //list of inovation numbers this network has (SORTED)
	id int //network id
	learningRate float64 //learning rate for backprop
	output []*Node //output nodes
	input []*Node //input nodes
	fitness float64
	adjustedFitness float64
	numInnovation int
	networkId int
}

//processes the network
func (n *Network) Process(input []float64) {
	for i := 0; i < len(n.input); i++ {
		n.input[i].setValue(input[i])
	}
}

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

func (n *Network) mutateConnection(from int, to int, innovation int) {
	//n.getNode(to).addRecCon(n.getNode(from).addSendCon(GetConnectionInstance(n.getNode(to), n.getNode(from), innovation)))

	c := n.getNode(from).addSendCon(GetConnectionInstance(n.getNode(from), n.getNode(to), innovation))
	b := n.getNode(to).addRecCon(c)

	fmt.Println("/////")
	fmt.Println(c)
	fmt.Println(b)
	fmt.Println("/////")
	n.addInnovation(innovation)

	n.getNode(to).numConIn++
	n.getNode(from).numConOut++
	n.numConnections++
}

func (n *Network) addInnovation(num int) {
	if len(n.innovation) <= n.numInnovation+1 {
		n.innovation = append(n.innovation,  num)
	} else {
		n.innovation[n.numInnovation+1] =  num
	}
	n.numInnovation++
}

//searches to remove the inovation
func (n *Network) removeInnovation(num int) {
	for i := 0; i < len(n.innovation); i++ {
		if n.innovation[i] == num {
			n.innovation = append(n.innovation[:i],n.innovation[i+1:]...)
		}
	}

	n.numInnovation--
}

/*
change from nodes connection to one with new node
change to nodes pointer to one sent by by new node
 */
func (n *Network) mutateNode(from int, to int, innovatonA int, innovationB int) int {
	fromNode := n.getNode(from)
	toNode := n.getNode(to)
	newNode := n.createNode()

	n.addInnovation(innovatonA)
	n.addInnovation(innovationB)

	//creates and modfies the connection to the toNode
	for i := 0; i < len(toNode.receive); i++ {
		if toNode.receive[i] != nil {
			fmt.Println(toNode.receive[i].nodeFrom.id)
		}
		if toNode.receive[i] != nil && fromNode == toNode.receive[i].nodeFrom { //compares the memory location
			toNode.receive[i] = newNode.addSendCon(GetConnectionInstance(newNode, toNode, innovatonA))
		}
	}
	fmt.Println("")
	//todo find a better way?
	for i := 0; i < len(fromNode.send); i++ {
		if fromNode.send[i].nodeTo != nil {
			fmt.Println(fromNode.send[i].nodeTo.id)
		}
		if fromNode.send[i].nodeTo != nil && fromNode.send[i].nodeTo.id == toNode.id {
			fromNode.send[i].nodeTo = newNode

			n.removeInnovation(fromNode.send[i].inNumber)
			fromNode.send[i].inNumber = innovationB

			newNode.addRecCon(&fromNode.send[i])
		}
	}

	fmt.Print(newNode)
	return newNode.id
}

func (n *Network) createNode() *Node {
	node := Node {value:0, numConOut: 0, numConIn: 0, influenceRecieved: 0, inputRecieved: 0, id:n.id, receive:make([]*Connection, len(n.input)), send:make([]Connection, len(n.output))}
	n.id++

	if (node.id+1) >= len(n.nodeList) {
		n.nodeList = append(n.nodeList,  node)
	} else {
		n.nodeList[len(n.nodeList)-(1+node.id)] =  node
	}

	return n.getNode(node.id)
}

func GetNetworkInstance(input int, output int, id int) Network {
	n := Network{numInnovation: 0, networkId: id, id: 0, learningRate: .1, numConnections:0, nodeList:make([]Node, (input+output)*2), output: make([]*Node, output), input: make([]*Node, input)}

	fmt.Print("initialized")

	//create output nodes
	for i := 0; i < output; i++ {
		n.output[i] = n.createNode()
	}

	fmt.Print("output")

	//creates the input nodes and adds them to the network
	for i := 0; i < input; i++ {
		n.input[i] = n.createNode()
		for a := 0; a < output; a++ {
			n.mutateConnection(n.input[i].id, n.output[a].id, n.numConnections)
		}
	}

	return n
}

func (n *Network) getNode(i int) *Node {
	return &n.nodeList[len(n.nodeList)-i-1]
}

