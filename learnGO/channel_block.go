package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch1 := make(chan int)
	go pump(ch1)       // pump hangs
	fmt.Println(<-ch1) // prints only 0

	c := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		c <- `foo` //unbuffered chan,
	}()

	go func() {
		defer wg.Done()

		time.Sleep(time.Second * 1)
		println(`Message: ` + <-c) // only receiver and sender both ready
	}()

	wg.Wait()
}

func pump(ch chan int) {
	for i := 0; ; i++ {
		ch <- i
	}
}
