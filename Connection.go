package main

import (
	"math/rand"
)
type Connection struct{
	weight float64
	disable bool
	nextWeight float64
	nodeTo *Node
	nodeFrom *Node
	inNumber int
}

func isRealConnection(c *Connection) bool {
	if c.nodeFrom == nil {
		return false
	}

	return true
}
//these act as the middle man between nodes
func (c *Connection) notifyValue() {
	if c.nodeTo != nil {
		c.nodeTo.recieveValue()
	}
}

func (c *Connection) notifyInfluence() {
	if c.nodeFrom != nil {
		c.nodeFrom.recieveInfluence()
	}
}

func GetConnectionInstance(from *Node, to *Node, inNumber int) Connection{
	return Connection{weight: rand.Float64()*.2 + .4, disable: false, nextWeight: 0, nodeTo: to, nodeFrom: from, inNumber: inNumber}
}