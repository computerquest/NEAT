package main

import ("fmt")

func main() {
	data := [][][]float64{
		{
			{},
			{},
		},
		{
			{},
			{},
		},
		{
			{},
			{},
		},
		{
			{},
			{},
		},
		{
			{},
			{},
		},
	}
	s := GetNeatInstance(40, 5, 5)

	s.start(data)

	fmt.Print("works")
}
