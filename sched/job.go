package sched

import "errors"

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

func (s *Scheduler) CreateNewJob(jobName string, args interface{}) (<-chan interface{}, error) {
	if jd, ok := s.jobDefs[jobName]; !ok {
		return nil, errors.New("JobDefinition not found")
	} else {
		var c chan interface{}
		if jd.hasRV {
			c = make(chan interface{}, 1)
		}
		j := Job{
			ref:          jd,
			args:         args,
			ceilPriority: jd.priority,
			resChan:      c,
		}
		if jd.preemptable {
			j.earlyBreak = make(chan interface{}, 1)
		}
		s.recvChan <- j
		return c, nil
	}
}
