package sched

import (
	"sync"
)

var priorityGap int = 255

type Captin struct {
	schedulers map[string]*Scheduler
}

type JdHandler func(chan struct{}, interface{}) interface{}

type JobDef struct {
	jobName      string
	preemptable  bool
	ceilPriority int // default is equal than `priority`
	priority     int
	handler      JdHandler
}

type Job struct {
	ref        *JobDef
	args       interface{}
	earlyBreak chan struct{}
	assoJob    *Job // sub-job in this job, it maybe under the different worker-pool
}

type Scheduler struct {
	schedName string
	workerNum int
	// pre-define
	jobMap  map[int][]*Job
	jobDefs map[string]*JobDef
	// runtime
	prioMap  map[int]int
	recvChan chan Job
	jobQueue chan Job
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

type Mutex struct {
	mu     sync.Mutex
	belong *Job
}
