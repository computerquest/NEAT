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
	learningRate    float64 //learning rate for backprop
	output          []*Node //output nodes
	input           []*Node //input nodes
	fitness         float64
	adjustedFitness float64
	networkId       int
	species         int
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
	n.input[input] = n.createNode(100) //starts unconnected and will form connections over time

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
func (n *Network) BackProp(input []float64, desired []float64) float64 {
	n.Process(input) //need to do so that you are perfkorming the algorithm on that set of values

	error := 0.0

	//this will calc all the influence
	for i := 0; i < len(n.output); i++ {
		n.output[i].setInfluence(n.output[i].value - desired[i])
		error += math.Abs(n.output[i].value - desired[i])
	}

	//actually adjusts the weights
	for i := 0; i < len(n.nodeList); i++ {
		//derivative := sigmoidDerivative(n.nodeList[i].value)
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
func (n *Network) trainSet(input [][][]float64, lim int) float64 {
	errorChange := -1000.0 //will be percent of error
	lastError := 1000.0

	bestWeight := make([][]float64, len(n.nodeList))
	for i := 0; i < len(n.nodeList); i++ {
		bestWeight[i] = make([]float64, len(n.nodeList[i].send))
	}

	n.resetWeight()

	strikes := 10
	for z := 1; strikes > 0 && z < lim && lastError > .000001; z++ {
		currentError := 0.0
		//resets all the next weights
		for i := 0; i < len(n.nodeList); i++ {
			for a := 0; a < len(n.nodeList[i].send); a++ {
				n.nodeList[i].send[a].nextWeight = 0
			}
		}

		for i := 0; i < len(input); i++ {
			currentError += n.BackProp(input[i][0], input[i][1])
		}

		//updates all the
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

	for i := 0; i < len(bestWeight); i++ {
		for a := 0; a < len(bestWeight[i]); a++ {
			n.nodeList[i].send[a].weight = bestWeight[i][a]
		}
	}

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
		n.innovation = n.innovation[0: len(n.innovation)+1]
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
	n.getNode(to).addRecCon(n.getNode(from).addSendCon(GetConnectionInstance(n.getNode(from), n.getNode(to), innovation)))
	n.addInnovation(innovation)
}
func (n *Network) numConnection() int {
	ans := 0
	for i := 0; i < len(n.nodeList); i++ {
		ans += len(n.nodeList[i].send)
	}

	return ans
}
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
func (n *Network) createNode(send int) *Node {
	node := Node{value: 0, influenceRecieved: 0, inputRecieved: 0, id: len(n.nodeList), receive: make([]*Connection, 0, 0), send: make([]Connection, 0, send)}

	if len(n.nodeList) >= cap(n.nodeList) {
		n.nodeList = append(n.nodeList, node)
	} else {
		n.nodeList = n.nodeList[0: len(n.nodeList)+1]
		n.nodeList[len(n.nodeList)-1] = node
	}

	return &n.nodeList[len(n.nodeList)-1]
}
func (n *Network) getNextNodeId() int {
	return len(n.nodeList)
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
			n.removeInnovation(toNode.receive[i].innovation)
			c := GetConnectionInstance(newNode, toNode, innovationB)
			c.weight = 1
			toNode.receive[i] = newNode.addSendCon(c)
		}
	}
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
func (network *Network) checkCircleMaster(n *Node, goal int) bool {
	preCheck := make([]int, len(network.nodeList))

	for i := 0; i < len(preCheck); i++ {
		preCheck[i] = i
	}

	return checkCircle(n, goal, preCheck)
}
func checkCircle(n *Node, goal int, preCheck []int) bool {
	ans := false
	if n.id == goal {
		return true
	}

	if preCheck[n.id] == -1 {
		return false
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
func clone(n *Network, in *[][]int) Network {
	a := GetNetworkInstance(len(n.input)-1, len(n.output), n.networkId, n.species, n.learningRate, false)

	for i := 0; i < len(n.nodeList)-len(n.input)-len(n.output); i++ {
		a.createNode(100)
	}

	for i := 0; i < len(n.innovation); i++ {
		a.mutateConnection((*in)[n.getInovation(i)][0], (*in)[n.getInovation(i)][1], n.getInovation(i))
	}

	for i := 0; i < len(a.nodeList); i++ {
		for b := 0; b < len(a.nodeList[i].send); b++ {
			if a.nodeList[i].send[b].innovation == n.nodeList[i].send[b].innovation {
				a.nodeList[i].send[b].weight = n.nodeList[i].send[b].weight
			} else {
				for c := 0; c < len(a.nodeList[i].send); c++ {
					if a.nodeList[i].send[b].innovation == n.nodeList[i].send[c].innovation {
						a.nodeList[i].send[b].weight = n.nodeList[i].send[c].weight
					}
				}
			}
		}
	}

	a.fitness = n.fitness

	return a
}
