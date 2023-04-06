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
	ceilPriority int  // default is equal than `priority`
	assoJob      *Job // sub-job in this job, it maybe under the different worker-pool
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
	jobQueue chan Job
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

type Mutex struct {
	mu           sync.Mutex
	resourceName string
	belong       *Job
}
