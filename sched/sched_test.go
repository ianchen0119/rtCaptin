package sched

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestFunctionality(t *testing.T) {
	c := NewCaptin()
	s, err := c.NewScheduler("http handler", 10)
	assert.EqualValues(t, nil, err, "NewScheduler() should not return error")
	// define jobs
	s.DefineNewJob("sleep", false, 0, true, sleep)
	ctx, cancel := context.WithCancel(context.Background())
	go s.Start(ctx)
	s.CreateNewJob("sleep", nil)
	time.Sleep(3 * time.Second)
	cancel()
}
