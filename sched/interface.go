package sched

import "context"

type captinOps interface {
	NewScheduler(string, int) (*Scheduler, error)
}

type schedulerOps interface {
	DefineNewTask(string, bool, int, bool,
		func(JobContext, interface{})) error
	CreateNewJob(string, interface{},
		[]*Resource) (<-chan interface{}, error)
	Start(context.Context)
}

type jobOps interface {
	NewJobContext() JobContext
}

// compile check
var _ captinOps = &Captin{}
var _ schedulerOps = &Scheduler{}
var _ jobOps = &Job{}
