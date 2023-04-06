package sched

import (
	"sync"
)

var priorityGap int = 255

type Captin struct {
	schedulers map[string]*Scheduler
}

type JdHandler func(chan interface{}, chan interface{}, interface{})

type JobDef struct {
	jobName     string
	preemptable bool
	priority    int
	hasRV       bool // return value of job's handler
	handler     JdHandler
}

type Job struct {
	ref          *JobDef
	args         interface{}
	earlyBreak   chan interface{}
	resChan      chan interface{}
	resources    []*Resource
	done         bool
	ceilPriority int // default is equal than `priority`
}

type Scheduler struct {
	schedName string
	workerNum int
	// pre-define
	jobMap  []*Job
	jobDefs map[string]*JobDef
	// runtime
	prioMap  map[int]int
	recvChan chan Job
	jobQueue chan *Job
	wg       sync.WaitGroup
}

type Resource struct {
	resourceName string
	res          interface{}
	belong       *Job
}
