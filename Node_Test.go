package main

import (
	"testing"
)

func TestSendCon(t *testing.T) {
	n := Node{send:make([]Connection, 3), receive: make([]*Connection, 3), numConIn: 0, numConOut: 0}

	n.addSendCon(Connection{nodeFrom:&n})

	if n.getSendCon(0).nodeFrom != &n || n.getSendCon(1).nodeFrom == &n {
		t.Errorf("Node the correct node")
	}
}
func TestRecCon(t *testing.T) {
	n := Node{send:make([]Connection, 3), receive: make([]*Connection, 3), numConIn: 0, numConOut: 0}

	n.addRecCon(&Connection{nodeTo: &n})
	if n.getRecCon(0).nodeTo != &n {
		t.Errorf("Node the correct node")
	}
}
