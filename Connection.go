package main

import (
	"math/rand"
)

type Connection struct {
	weight     float64 //weight of the connection
	disable    bool    //is the connection active
	nextWeight float64 //sum of backpropogation changes to weight
	nodeTo     *Node   //receiving node
	nodeFrom   *Node   //sending node
	innovation int     //innovation number
}

//sets the connection to a random weight
func (c *Connection) randWeight() {
	c.weight = rand.Float64()*.2 + .4
}

//notify nodeTo of nodeFrom value
func (c *Connection) notifyValue() {
	//if c.nodeTo != nil {
		c.nodeTo.recieveValue()
	//}
}

//notify nodeFrom of nodeTo influence
func (c *Connection) notifyInfluence() {
	//if c.nodeFrom != nil {
		c.nodeFrom.recieveInfluence()
	//}
}

//initialize connection
func GetConnectionInstance(from *Node, to *Node, inNumber int) Connection {
	c := Connection{weight: 0, disable: false, nextWeight: 0, nodeTo: to, nodeFrom: from, innovation: inNumber}
	c.randWeight()

	return c
}
