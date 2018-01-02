package main

import (
	"fmt"
	"math"
)

//100 NODE MAX!!!!!!!!!!!!!!!
//NOTE most of the calculating work is networked by nodes inside the struct
type Network struct {
	nodeList        []Node  //master list of nodes
	innovation      []int   //list of innovation numbers this network has (innovation number = unique connection)
	learningRate    float64 //learning rate for backprop
	output          []*Node //output nodes
	input           []*Node //input nodes
	fitness         float64 //effectiveness of network (1/error)
	adjustedFitness float64 //fitness relative to that in species
	networkId       int
	species         int //id of species
}

func GetNetworkInstance(input int, output int, id int, species int, learningRate float64, addCon bool) Network {
	n := Network{networkId: id, learningRate: learningRate, nodeList: make([]Node, 0, 100), output: make([]*Node, output), input: make([]*Node, input+1), species: species}

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
	n.input[input] = n.createNode(100) //bias starts unconnected and will form connections over time

	return n
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
//evaluate input returns values of output
func (n *Network) Process(input []float64) []float64 {
	//set input values
	for i := 0; i < len(n.input); i++ {
		if i < len(input) {
			n.input[i].setValue(input[i])
		} else {
			n.input[i].setValue(1)
		}
	}

	//values are calculated via connections and nodes signalling

	ans := make([]float64, len(n.output))
	for i := 0; i < len(n.output); i++ {
		ans[i] = n.output[i].value
	}

	return ans
}

//trains network for input values against the desired values and returns the error
func (n *Network) BackProp(input []float64, desired []float64) float64 {
	n.Process(input) //set the values for the input

	error := 0.0 //return value

	//this will calc all the influence
	for i := 0; i < len(n.output); i++ {
		n.output[i].setInfluence(n.output[i].value - desired[i])
		error += math.Abs(n.output[i].value - desired[i])
	}

	//all the influence is set the same way as values so it is set via connections and signalling

	//actually adjusts the weights
	for i := 0; i < len(n.nodeList); i++ {
		for a := 0; a < len(n.nodeList[i].receive); a++ {
			if n.nodeList[i].receive[a] != nil {
				if n.nodeList[i].receive[a].disable {
					continue
				}
				n.nodeList[i].receive[a].nextWeight -= (n.nodeList[i].receive[a].nodeFrom.value) * n.nodeList[i].influence * n.learningRate
			}
		}
	}

	return error
}

//the handling function for training returns the fitness
func (n *Network) trainSet(input [][][]float64, lim int) float64 {
	errorChange := -1000.0 //percent of error change
	lastError := 1000.0

	//initializes best weights
	bestWeight := make([][]float64, len(n.nodeList))
	for i := 0; i < len(n.nodeList); i++ {
		bestWeight[i] = make([]float64, len(n.nodeList[i].send))
	}

	n.resetWeight() //clears the current weight values

	strikes := 10 //number of times in a row that error can increase
	for z := 1; strikes > 0 && z < lim && lastError > .000001; z++ {
		currentError := 0.0

		//resets all the nextWeights
		for i := 0; i < len(n.nodeList); i++ {
			for a := 0; a < len(n.nodeList[i].send); a++ {
				n.nodeList[i].send[a].nextWeight = 0
			}
		}

		//trains each input
		for i := 0; i < len(input); i++ {
			currentError += n.BackProp(input[i][0], input[i][1])
		}

		//updates all the weight
		for i := 0; i < len(n.nodeList); i++ {
			for a := 0; a < len(n.nodeList[i].send); a++ {
				if n.nodeList[i].send[a].disable {
					continue
				}
				n.nodeList[i].send[a].weight += n.nodeList[i].send[a].nextWeight / float64(len(input))
			}
		}

		errorChange = (currentError - lastError) / lastError
		lastError = currentError

		//decreases the number of strikes or resets them and changes best weight
		if errorChange >= 0 {
			strikes--
		} else {
			for i := 0; i < len(n.nodeList); i++ {
				for a := 0; a < len(n.nodeList[i].send); a++ {
					bestWeight[i][a] = n.nodeList[i].send[a].weight
				}
			}
			strikes = 10
		}
	}

	//sets the weights back to the best
	for i := 0; i < len(bestWeight); i++ {
		for a := 0; a < len(bestWeight[i]); a++ {
			n.nodeList[i].send[a].weight = bestWeight[i][a]
		}
	}

	//calculate the final error
	final := 0.0
	for i := 0; i < len(input); i++ {
		stuff := n.Process(input[i][0])
		for a := 0; a < len(stuff); a++ {
			final += math.Abs(stuff[a] - input[i][1][a])
		}
	}

	n.fitness = 1 / final
	return final
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
//adds a connection from from (node id) to to (node id) and adds innovation to the network
func (n *Network) mutateConnection(from int, to int, innovation int) {
	n.getNode(to).addRecCon(n.getNode(from).addSendCon(GetConnectionInstance(n.getNode(from), n.getNode(to), innovation)))
	n.addInnovation(innovation)
}

//returns the number of connections
func (n *Network) numConnection() int {
	ans := 0
	for i := 0; i < len(n.nodeList); i++ {
		ans += len(n.nodeList[i].send)
	}

	return ans
}

//resets a networks weights
func (n *Network) resetWeight() {
	for i := 0; i < len(n.nodeList); i++ {
		for a := 0; a < len(n.nodeList[i].send); a++ {
			n.nodeList[i].send[a].randWeight()
			n.nodeList[i].send[a].nextWeight = 0
		}
	}
}

///////////////////////////////////////////////////////NODE
func (n *Network) getNode(i int) *Node {
	return &n.nodeList[i]
}

//returns a pointer to node created with a cap of send (parameter) for send field in node (already added to nodeList)
func (n *Network) createNode(send int) *Node {
	node := Node{value: 0, influenceRecieved: 0, inputRecieved: 0, id: len(n.nodeList), receive: make([]*Connection, 0, 0), send: make([]Connection, 0, send)}

	if len(n.nodeList) >= cap(n.nodeList) {
		n.nodeList = append(n.nodeList, node)
	} else {
		n.nodeList = n.nodeList[0 : len(n.nodeList)+1]
		n.nodeList[len(n.nodeList)-1] = node
	}

	return &n.nodeList[len(n.nodeList)-1]
}

//provides the id for the next node added
func (n *Network) getNextNodeId() int {
	return len(n.nodeList)
}

//adds a node in between from and to and adds the innovation numbers innovationA and innovationB and removes the original innovation number and returns the new node's id
func (n *Network) mutateNode(from int, to int, innovationA int, innovationB int) int {

	fromNode := n.getNode(from)
	toNode := n.getNode(to)
	newNode := n.createNode(100)

	n.addInnovation(innovationA)
	n.addInnovation(innovationB)

	//changes the connection recieved by toNode to a connection sent by newNode
	for i := 0; i < len(toNode.receive); i++ {
		if toNode.receive[i] != nil && fromNode == toNode.receive[i].nodeFrom {
			n.removeInnovation(toNode.receive[i].innovation)
			c := GetConnectionInstance(newNode, toNode, innovationB)
			toNode.receive[i] = newNode.addSendCon(c)
		}
	}

	//modifies the connection from fromNode by changing the toNode for the connection to newNode from toNode
	for i := 0; i < len(fromNode.send); i++ {
		if fromNode.send[i].nodeTo != nil && fromNode.send[i].nodeTo.id == toNode.id {
			fromNode.send[i].nodeTo = newNode
			fromNode.send[i].innovation = innovationA

			newNode.addRecCon(&fromNode.send[i])
		}
	}

	return newNode.id
}

///////////////////////////////////////////////////////MISC
//the handler for recursive checkCircle makes sure that there is no circular structure in the network returns if from node n there is some form of connection between node goal (id)
func (network *Network) checkCircleMaster(n *Node, goal int) bool {
	preCheck := make([]int, len(network.nodeList))

	for i := 0; i < len(preCheck); i++ {
		preCheck[i] = i
	}

	return checkCircle(n, goal, preCheck)
}

//takes the node to check from, the target node, and a precheck to see if a node has already been checked for connection
func checkCircle(n *Node, goal int, preCheck []int) bool {
	ans := false
	if n.id == goal {
		return true
	}

	//checks for the precheck
	if preCheck[n.id] == -1 {
		return false
	}

	//checks next stop down
	for i := 0; i < len(n.receive); i++ {
		ans = checkCircle(n.receive[i].nodeFrom, goal, preCheck)
		if ans {
			break
		}
	}

	//sets the precheck
	if !ans {
		preCheck[n.id] = -1
	}

	return ans
}

//takes n (network to clone) and in (master innovation list) and returns a duplicate of the network
func clone(n *Network) Network {

	//need to totally reconstruct because otherwise the pointers in connections and such would be screwed up
	ans := GetNetworkInstance(len(n.input)-1, len(n.output), n.networkId, n.species, n.learningRate, false)

	for i := 0; i < len(n.nodeList)-len(n.input)-len(n.output); i++ {
		ans.createNode(100)
	}

	for i := 0; i < len(n.nodeList); i++ {
		for a := 0; a < len(n.nodeList[i].send); a++ {
			ans.mutateConnection(n.nodeList[i].send[a].nodeFrom.id, n.nodeList[i].send[a].nodeTo.id, n.nodeList[i].send[a].innovation)
			ans.nodeList[i].send[a].weight = n.nodeList[i].send[a].weight
		}
	}

	ans.fitness = n.fitness

	return ans
}
