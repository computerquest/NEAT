package main

type Node struct {
	value float64
	id int
	receive []*Connection
	send []Connection //this list is seqential for initialization
}
