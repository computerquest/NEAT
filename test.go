package main

import (
	"fmt"
	"sync"
)

func main() {
	/*c := increment()
	cSum := puller(c)
	for n := range cSum {
		fmt.Print(n)
		fmt.Println("done")
	}*/

	var w sync.WaitGroup

	fmt.Println("starting")
	ch := make(chan int)

	w.Add(2)

	go func() {
		//fmt.Println("0")
		for i := 0; i < 10; i++ {
			fmt.Println("writing")
			ch <- i
			ch <- (i + 1)
		}

		close(ch)
		w.Done()
	}()

	go func() {
		//fmt.Println("1")
		for n := range ch {
			fmt.Println("reading: ", n)
		}
		w.Done()
	}()

	w.Wait()
}

func increment() chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Print(i)
			fmt.Println("here")

			out <- i
		}

		fmt.Println("closing increment")
		close(out)
	}()

	fmt.Println("exiting increment")
	return out
}

func puller(c chan int) chan int {
	out := make(chan int)

	go func() {
		var sum int
		for n := range c {
			fmt.Print(n)
			fmt.Println("there")
			sum += n
		}
		out <- sum
		fmt.Println("closing puller")

		close(out)
	}()

	fmt.Println("exiting puller")

	return out
}
