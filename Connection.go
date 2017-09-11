package main

type Connection struct{
	weight float64
	//sendValue *float64
	disable bool
	nextWeight float64
	//connectInfluence *float64
	idTo int
	idFrom int
}
