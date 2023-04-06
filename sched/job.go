package sched

import (
	"errors"
)

func (s *Scheduler) DefineNewJob(jobName string, preemptable bool,
	priority int, hasReturnVal bool,
	handler func(chan interface{}, chan interface{}, interface{})) error {
	if priority > priorityGap {
		return errors.New("Job's priority should not greater than 255")
	}

	jd := &JobDef{
		jobName:     jobName,
		preemptable: preemptable,
		priority:    priority,
		handler:     handler,
		hasRV:       hasReturnVal,
	}

	s.jobDefs[jobName] = jd

	return nil
}

func (s *Scheduler) CreateNewJob(jobName string, args interface{}, resources []*Resource) (<-chan interface{}, error) {
	if jd, ok := s.jobDefs[jobName]; !ok {
		return nil, errors.New("JobDefinition not found")
	} else {
		var c chan interface{}
		if jd.hasRV {
			c = make(chan interface{}, 1)
		}
		if args == nil {
			args = struct{}{}
		}
		j := Job{
			ref:          jd,
			args:         args,
			ceilPriority: jd.priority,
			resChan:      c,
			done:         false,
			resources:    resources,
		}
		if jd.preemptable {
			j.earlyBreak = make(chan interface{}, 1)
		}
		s.recvChan <- j
		return c, nil
	}
}

func NewResource(resName string, val interface{}) *Resource {
	return &Resource{
		resourceName: resName,
		res:          val,
	}
}
