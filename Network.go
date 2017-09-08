package main

import "fmt"

type Network struct {
	nodes [][]Node
	numConnections int
	numNodes int
	innovation []int
	id int
}

func (n *Network) process(input []float64) {
	for i := 0; i < len(n.nodes[0]); i++ {
		n.nodes[0][i].value = input[i];
	}

	for i := 1; i < len(n.nodes); i++ {
		for a := 0; a < len(n.nodes[i]); a++ {
			n.nodes[i][a].value = sigmoid(n.nodes[i][a].netInput())
		}
	}
}

//todo needs to be finished
func sigmoid(value float64) float64 {
	return 0;
}

//todo need to be tested
func (n *Network) GetInstance(input int, output int) {
	n.numConnections = 0
	n.numNodes = 0

	n.nodes = make([][]Node, 3)

	n.nodes[0] = make([]Node, input);
	n.nodes[2] = make([]Node, output);

	fmt.Print("initialized")

	//creates the nodes and adds them to the network
	for i := 0; i < input; i++ {
		n.nodes[0][i] = Node {value:0, id:n.id, send:make([]Connection, output)}

		//creates the connections
		for a := 0; a < output; a++ {
			n.nodes[0][i].send[a] = Connection{weight: 1, disable:false, sendValue:&n.nodes[0][i].value}
			n.numConnections++;
		}

		n.id++;
	}

	fmt.Print("input")

	for i := 0; i < output; i++ {
		n.nodes[2][i] = Node {value:0, id:n.id, receive:make([]*Connection, input)}
		n.id++;
	}

	fmt.Print("output")

	//populates output recieve
	for i := 0; i < output; i++ {
		for a := 0; a < input; a++ {
			n.nodes[2][i].receive[a] = &n.nodes[0][a].send[i]
		}
	}

	fmt.Print("recieving connection")

	n.numNodes = input+output;
}
