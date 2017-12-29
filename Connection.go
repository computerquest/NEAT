package main

import (
	"math/rand"
)

type Connection struct {
	weight     float64
	disable    bool
	nextWeight float64
	nodeTo     *Node
	nodeFrom   *Node
	innovation int
}

func (c *Connection) randWeight() {
	c.weight = rand.Float64()*.2 + .4
}

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

func GetConnectionInstance(from *Node, to *Node, inNumber int) Connection {
	c := Connection{weight: 0, disable: false, nextWeight: 0, nodeTo: to, nodeFrom: from, innovation: inNumber}
	c.randWeight()

	return c
}
