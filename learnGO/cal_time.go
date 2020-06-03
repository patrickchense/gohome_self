package main

import (
	"flag"
	"gohome_self/learnGO/log"
	"time"
)

func main() {

	tc := flag.Int("threads", 10, "Thread Count")
	rut := flag.Int("rampup", 30, "Ramp up time in seconds")
	et := flag.Int("etime", 1, "Execution time in minutes")
	flag.Parse()

	//Check if execution time is more than ramp up time
	if *et*60 < *rut {
		log.Fatalln("Total execution time needs to be more than ramp up time")
	}

	waitTime := *rut / *tc

	log.Printf("Execution will happen with %d users with a ramp up time of %d seconds for %d minutes\n", *tc, *rut, *et)

	tchan := make(chan int)
	go func(c chan<- int) {
		for ti := 1; ti <= *tc; ti++ {
			log.Printf("Thread Count %d", ti)
			c <- ti
			time.Sleep(time.Duration(waitTime) * time.Second)
		}
	}(tchan)

	timeout := time.After(time.Duration(*et*60) * time.Second)

	for {
		select { //Select blocks the flow until one of the channels receives a message
		case <-timeout: //receives a msg when execution duration is over
			log.Printf("Execution completed")
			return
		case ts := <-tchan: //receives a message when a new user thread has to be initiated
			log.Printf("Thread No %d started", ts)
			go func(t int) {
				//This is the place where you add all your tests
				//In my case they were making rpc calls over rabbitmq with random inputs
				//They keep running till the end of execution
				for {
					//sample test
					test()
				}
			}(ts)
		}
	}
}
