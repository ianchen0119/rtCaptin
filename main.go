package main

import (
	"context"
	"fmt"
	"time"

	cap "github.com/ianchen0119/rtCaptin/sched"
)

func sleep(notify chan struct{}, args interface{}) interface{} {
	for {
		select {
		case <-notify:
			return nil
		default:
			fmt.Println("Good morning!")
			return struct{}{}
		}
	}
}

func sayHi(notify chan struct{}, args interface{}) interface{} {
	for {
		select {
		case <-notify:
			return nil
		default:
			fmt.Println("Hi!")
			return struct{}{}
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
	s.DefineNewJob("sleep", false, 0, sleep)
	s.DefineNewJob("hi", false, 255, sayHi)
	ctx, cancel := context.WithCancel(context.Background())
	go s.Start(ctx)
	for i := 0; i < 1000; i++ {
		s.CreateNewJob("sleep", nil)
		if i == 999 {
			s.CreateNewJob("hi", nil)
		}
	}
	time.Sleep(10 * time.Second)
	cancel()
}
