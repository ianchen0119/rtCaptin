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
			time.Sleep(1 * time.Second)
			fmt.Println("Good morning!")
			return struct{}{}
		}
	}
}

func main() {
	c := cap.NewCaptin()
	s, err := c.NewScheduler("http handler", 10)
	if err != nil {
		return
	}
	// define jobs
	s.DefineNewJob("sleep", false, 0, sleep)
	ctx, cancel := context.WithCancel(context.Background())
	go s.Start(ctx)
	s.CreateNewJob("sleep", nil)
	s.CreateNewJob("sleep", nil)
	s.CreateNewJob("sleep", nil)
	s.CreateNewJob("sleep", nil)
	s.CreateNewJob("sleep", nil)
	s.CreateNewJob("sleep", nil)
	s.CreateNewJob("sleep", nil)
	s.CreateNewJob("sleep", nil)
	s.CreateNewJob("sleep", nil)
	s.CreateNewJob("sleep", nil)
	s.CreateNewJob("sleep", nil)
	s.CreateNewJob("sleep", nil)
	time.Sleep(10 * time.Second)
	cancel()
}
