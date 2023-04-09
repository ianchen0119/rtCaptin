package sched

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func sleep(ctx JobContext, args interface{}) {
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

func TestFunctionality(t *testing.T) {
	c := NewCaptin()
	s, err := c.NewScheduler("http handler", 10)
	assert.EqualValues(t, nil, err, "NewScheduler() should not return error")
	// define jobs
	s.DefineNewTask("sleep", false, 0, true, sleep)
	ctx, cancel := context.WithCancel(context.Background())
	go s.Start(ctx)
	s.CreateNewJob("sleep", nil, nil)
	time.Sleep(3 * time.Second)
	cancel()
}
