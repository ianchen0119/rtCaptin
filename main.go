package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	cap "github.com/ianchen0119/rtCaptin/sched"
)

func sleep(ctx cap.JobContext, args interface{}) {
	notify := ctx.GetEarlyBreakChan()
	res := ctx.GetResultChan()
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

func sayHi(ctx cap.JobContext, args interface{}) {
	notify := ctx.GetEarlyBreakChan()
	for {
		select {
		case <-notify:
			return
		default:
			time.Sleep(3 * time.Millisecond)
			fmt.Println("Hi!")
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
	// define tasks
	s.DefineNewTask("sleep", false, 0, true, sleep)
	s.DefineNewTask("hi", false, 20, false, sayHi)
	ctx, cancel := context.WithCancel(context.Background())
	go s.Start(ctx)
	res, _ := s.CreateNewJob("sleep", "Good Morning!\n", nil)
	var mu sync.Mutex
	rc := cap.NewResource("lock", mu)
	rcList := []*cap.Resource{}
	rcList = append(rcList, rc)
	for i := 0; i < 1000; i++ {
		s.CreateNewJob("sleep", "Good Morning!\n", nil)
		if i == 100 {
			s.CreateNewJob("hi", nil, rcList)
		}
		if i == 102 {
			s.CreateNewJob("sleep", "Good Morning!\n", rcList)
		}
	}
	val := <-res
	fmt.Println(val)
	time.Sleep(10 * time.Second)
	cancel()
}
