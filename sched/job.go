package sched

import "errors"

func (s *Scheduler) DefineNewJob(jobName string, preemptable bool,
	priority int, handler func(chan struct{}, interface{}) interface{}) error {
	if priority > priorityGap {
		return errors.New("Job's priority should not greater than 255")
	}

	jd := &JobDef{
		jobName:      jobName,
		preemptable:  preemptable,
		priority:     priority,
		ceilPriority: priority,
		handler:      handler,
	}

	s.jobDefs[jobName] = jd

	return nil
}

// TODO: get the result of the executed Job
func (s *Scheduler) CreateNewJob(jobName string, args interface{}) error {
	if jd, ok := s.jobDefs[jobName]; !ok {
		return errors.New("JobDefinition not found")
	} else {
		j := Job{
			ref:     jd,
			assoJob: nil,
			args:    args,
		}
		if jd.preemptable {
			j.earlyBreak = make(chan struct{}, 1)
		}
		s.recvChan <- j
		return nil
	}
}
