package main

import (
	"testing"
)

func TestAddSendCon(t *testing.T) {
	n := Node{send:make([]Connection, 3), receive: make([]*Connection, 3), numConIn: 0, numConOut: 0}

	n.addSendCon(Connection{nodeFrom:&n})
}
func TestAddRecCon(t *testing.T) {
	n := Node{send:make([]Connection, 3), receive: make([]*Connection, 3), numConIn: 0, numConOut: 0}

	n.addRecCon(&Connection{})
}

func TestGetSendCon(t *testing.T) {
	n := Node{send:make([]Connection, 3), receive: make([]*Connection, 3), numConIn: 0, numConOut: 0}

	for i := 0; i < len(n.send); i++ {

	}
}

func TestGetRecCon(t *testing.T) {

}
