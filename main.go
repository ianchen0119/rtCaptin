package main

import (
	"context"
	"fmt"
	"time"

	cap "github.com/ianchen0119/rtCaptin/sched"
)

func sleep(notify chan interface{}, res chan interface{}, args interface{}) {
	for {
		select {
		case <-notify:
			return
		default:
			time.Sleep(5 * time.Millisecond)
			fmt.Printf("%s", args)
			res <- 1
			return
		}
	}
}

func sayHi(notify chan interface{}, res chan interface{}, args interface{}) {
	for {
		select {
		case <-notify:
			return
		default:
			time.Sleep(3 * time.Millisecond)
			fmt.Println("Hi!")
			res <- 1
			return
		}
	}
}

func main() {
	c := cap.NewCaptin()
	s, err := c.NewScheduler("http handler", 1)
	if err != nil {
		return
	}
	// define jobs
	s.DefineNewJob("sleep", false, 0, true, sleep)
	s.DefineNewJob("hi", false, 255, false, sayHi)
	ctx, cancel := context.WithCancel(context.Background())
	go s.Start(ctx)
	res, _ := s.CreateNewJob("sleep", "Good Morning!\n")
	for i := 0; i < 100; i++ {
		s.CreateNewJob("sleep", "Good Morning!\n")
		if i == 60 {
			s.CreateNewJob("hi", nil)
		}
	}
	val := <-res
	fmt.Println(val)
	time.Sleep(10 * time.Second)
	cancel()
}
