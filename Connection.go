package main

type Connection struct{
	weight float64
	disable bool
	nextWeight float64
	nodeTo *Node
	nodeFrom *Node
	inNumber int
}

//these act as the middle man between nodes
func (c *Connection) notifyValue() {
	c.nodeTo.recieveValue()
}

func (c *Connection) notifyInfluence() {
	c.nodeFrom.recieveInfluence()
}

func GetConnectionInstance(from *Node, to *Node, inNumber int) Connection{
	return Connection{weight: 1, disable: false, nextWeight: 0, nodeTo: to, nodeFrom: from, inNumber: inNumber}
}