package main

type Node struct {
	value float64
	id int
	receive []*Connection
	send []Connection //this list is seqential for initialization
}

func (n Node) netInput() float64 {
	var sum float64 = 0
	for i := 0; i < len(n.receive); i++ {
		c := n.receive[i]
		if !c.disable {
			sum += (*c.sendValue)*c.weight
		}
	}

	return sum
}
