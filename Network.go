package main

import (
	"fmt"
	"math"
)

//100 NODE MAX!!!!!!!!!!!!!!!
//NOTE most of the calculating work is networked by nodes inside the struct
type Network struct {
	nodeList        []Node  //master list of nodes
	innovation      []int   //list of inovation numbers this network has (SORTED)
	id              int     //network id
	learningRate    float64 //learning rate for backprop
	output          []*Node //output nodes
	input           []*Node //input nodes
	fitness         float64
	adjustedFitness float64
	networkId       int
	species         int
}

func GetNetworkInstance(input int, output int, id int, species int, learningRate float64, addCon bool) Network {
	n := Network{networkId: id, id: 0, learningRate: learningRate, nodeList: make([]Node, 0, 100), output: make([]*Node, output), input: make([]*Node, input+1), species: species}

	//create output nodes
	for i := 0; i < output; i++ {
		n.output[i] = n.createNode(0)
	}

	//creates the input nodes and adds them to the network
	startInov := 0 //this should work
	for i := 0; i < input; i++ {
		n.input[i] = n.createNode(100)
		if addCon {
			for a := 0; a < output; a++ {
				n.mutateConnection(n.input[i].id, n.output[a].id, startInov)
				startInov++
			}
		}
	}
	n.input[input] = n.createNode(100) //starts unconnected and will form connections over time

	return n
}
func isRealNetwork(n *Network) bool {
	if cap(n.nodeList) != 0 {
		return true
	}

	return false
}

func printNetwork(n *Network) {
	fmt.Println("network id: ", n.networkId, " species id: ", n.species)
	fmt.Print("expected connection: ", n.innovation)
	fmt.Println()

	for b := 0; b < len(n.nodeList); b++ {
		fmt.Print("node: ", n.nodeList[b].id, " sending: ")
		for c := 0; c < len(n.nodeList[b].send); c++ {
			fmt.Print(n.nodeList[b].send[c].nodeTo.id, " ")
		}

		fmt.Print("receive: ")
		for c := 0; c < len(n.nodeList[b].receive); c++ {
			fmt.Print(n.nodeList[b].receive[c].nodeFrom.id, " ")
		}

		fmt.Println()
	}
}

////////////////////////////////////////////////////////////RUNNING
func (n *Network) Process(input []float64) []float64 {
	for i := 0; i < len(n.input); i++ {
		if i < len(input) {
			n.input[i].setValue(input[i])
		} else {
			n.input[i].setValue(1)
		}
	}

	ans := make([]float64, len(n.output))
	for i := 0; i < len(n.output); i++ {
		ans[i] = n.output[i].value
	}

	return ans
}

func (network *Network) checkCircleMaster(n *Node, goal int) bool {
	preCheck := make([]int, len(network.nodeList))

	for i := 0; i < len(preCheck); i++ {
		preCheck[i] = i
	}

	return checkCircle(n, goal, preCheck)
}

//TODO: fix method so it does not depend on all ids for nodes being in order
func checkCircle(n *Node, goal int, preCheck []int) bool {
	ans := false
	if n.id == goal {
		return true
	}

	if preCheck[n.id] == -1 {
		//return false
	}

	//checks for cirular dependency
	for i := 0; i < len(n.receive); i++ {
		ans = checkCircle(n.receive[i].nodeFrom, goal, preCheck)
		if ans {
			break
		}
	}

	if !ans {
		preCheck[n.id] = -1
	}

	return ans
}
func (n *Network) BackProp(input []float64, desired []float64) float64 {
	n.Process(input) //need to do so that you are perfkorming the algorithm on that set of values

	var error float64

	//this will calc all the influence
	for i := 0; i < len(n.output); i++ {
		n.output[i].setInfluence(n.output[i].value - desired[i])
		error += n.output[i].value - desired[i]
	}

	//actually adjusts the weights
	for i := 0; i < len(n.nodeList); i++ {
		derivative := sigmoidDerivative(n.nodeList[i].value)
		for a := 0; a < len(n.nodeList[i].receive); a++ {
			if n.nodeList[i].receive[a] != nil {
				if math.IsNaN(derivative) || math.IsNaN(n.nodeList[i].receive[a].nodeFrom.value) || math.IsNaN(n.nodeList[i].influence) || math.IsNaN(n.nodeList[i].receive[a].nextWeight) || math.IsNaN(n.learningRate) {
					fmt.Print(derivative, (n.nodeList[i].receive[a].nodeFrom.value), n.nodeList[i].influence, n.learningRate)
					printNetwork(n)
					fmt.Print("alert the master")
				}
				n.nodeList[i].receive[a].nextWeight += derivative * (n.nodeList[i].receive[a].nodeFrom.value) * n.nodeList[i].influence * n.learningRate
			}
		}
	}

	return error
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
func (n *Network) trainSet(input [][][]float64, lim int) float64 {
	errorChange := -1000.0 //will be percent of error

	lastError := 1000.0
	strikes := 10
	for z := 1; strikes > 0 && lastError > .00000000001 && z < lim; z++ {
		currentError := 0.0
		//resets all the next weights
		for i := 0; i < len(n.nodeList); i++ {
			for a := 0; a < len(n.nodeList[i].send); a++ {
				if isRealConnection(&n.nodeList[i].send[a]) {
					n.nodeList[i].send[a].nextWeight = 0
				}
			}
		}

		for i := 0; i < len(input); i++ {
			currentError += n.BackProp(input[i][0], input[i][1])
		}

		//updates all the
		for i := 0; i < len(n.nodeList); i++ {
			for a := 0; a < len(n.nodeList[i].send); a++ {
				if isRealConnection(&n.nodeList[i].send[a]) {
					n.nodeList[i].send[a].weight += n.nodeList[i].send[a].nextWeight / float64(len(input))
				}
			}
		}

		errorChange = (currentError - lastError) / lastError
		fmt.Printf("Gen: %d Current Error: %f avg: %e change: %e percent change: %f", z, currentError, currentError/float64(len(input)), currentError-lastError, errorChange)
		fmt.Println()
		lastError = currentError

		if errorChange > -.01 {
			//strikes--
		}
	}

	n.fitness = 1 / lastError //TODO: could be bad
	return lastError
}

/////////////////////////////////////////////////////////INNOVATION
func (n *Network) getInovation(pos int) int {
	return n.innovation[pos]
}
func (n *Network) addInnovation(num int) {
	if len(n.innovation) >= cap(n.innovation) {
		n.innovation = append(n.innovation, num)
	} else {
		n.innovation = n.innovation[0 : len(n.innovation)+1]
		n.innovation[len(n.innovation)-1] = num
	}
}
func (n *Network) containsInnovation(num int) bool {
	for i := 0; i < len(n.innovation); i++ {
		if n.innovation[i] == num {
			return true
		}
	}

	return false
}
func (n *Network) removeInnovation(num int) {
	for i := 0; i < len(n.innovation); i++ {
		if n.innovation[i] == num {
			n.innovation = append(n.innovation[:i], n.innovation[i+1:]...)
		}
	}
}

/////////////////////////////////////////////////////////CONNECTION
func (n *Network) mutateConnection(from int, to int, innovation int) {
	c := n.getNode(from).addSendCon(GetConnectionInstance(n.getNode(from), n.getNode(to), innovation))
	n.getNode(to).addRecCon(c)

	n.addInnovation(innovation)
}
func (n *Network) numConnection() int {
	ans := 0
	for i := 0; i < len(n.nodeList); i++ {
		ans += len(n.nodeList[i].send)
	}

	return ans
}

///////////////////////////////////////////////////////NODE
func (n *Network) getNode(i int) *Node {
	return &n.nodeList[i]
}
func (n *Network) createNode(send int) *Node {
	node := Node{value: 0, influenceRecieved: 0, inputRecieved: 0, id: n.id, receive: make([]*Connection, 0, 0), send: make([]Connection, 0, send)}
	n.id++

	if len(n.nodeList) >= cap(n.nodeList) {
		n.nodeList = append(n.nodeList, node)
	} else {
		n.nodeList = n.nodeList[0 : len(n.nodeList)+1]
		n.nodeList[len(n.nodeList)-1] = node
	}

	return &n.nodeList[len(n.nodeList)-1]
}
func (n *Network) getNextNodeId() int {
	return n.id
}
func (n *Network) mutateNode(from int, to int, innovationA int, innovationB int) int {
	fromNode := n.getNode(from)
	toNode := n.getNode(to)
	newNode := n.createNode(100)

	n.addInnovation(innovationA)
	n.addInnovation(innovationB)

	//creates and modfies the connection to the toNode
	for i := 0; i < len(toNode.receive); i++ {
		if toNode.receive[i] != nil && fromNode == toNode.receive[i].nodeFrom { //compares the memory location
			toNode.receive[i] = newNode.addSendCon(GetConnectionInstance(newNode, toNode, innovationB))
		}
	}
	//todo find a better way?
	for i := 0; i < len(fromNode.send); i++ {
		if fromNode.send[i].nodeTo != nil && fromNode.send[i].nodeTo.id == toNode.id {
			fromNode.send[i].nodeTo = newNode

			n.removeInnovation(fromNode.send[i].inNumber)
			fromNode.send[i].inNumber = innovationA

			newNode.addRecCon(&fromNode.send[i])
		}
	}

	return newNode.id
}
