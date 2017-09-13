package main

type Connection struct{
	weight float64
	disable bool
	nextWeight float64
	nodeTo *Node
	nodeFrom *Node
}

//these act as the middle man between nodes
func (c *Connection) notifyValue() {
	c.nodeTo.recieveValue()
}

func (c *Connection) notifyInfluence() {
	c.nodeFrom.recieveInfluence()
}