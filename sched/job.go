package sched

import (
	"context"
	"errors"
)

func (s *Scheduler) DefineNewJob(jobName string, preemptable bool,
	priority int, hasReturnVal bool,
	handler func(JobContext, interface{})) error {
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
			owner:        s,
		}
		if jd.preemptable {
			j.earlyBreak = make(chan interface{}, 1)
		}
		s.recvChan <- j
		return c, nil
	}
}

func (j *Job) NewJobContext() JobContext {
	ctxMap := JobContextMap{map[string]chan interface{}{
		"earlyBreak": j.earlyBreak,
		"result":     j.resChan,
	}}
	return JobContext{
		ctx: context.WithValue(context.Background(),
			"content", ctxMap),
	}
}

func (c *JobContext) GetEarlyBreakChan() chan interface{} {
	m := c.ctx.Value("content")
	if res, ok := m.(JobContextMap).m["earlyBreak"]; ok {
		return res
	}
	return nil
}

func (c *JobContext) GetResultChan() chan interface{} {
	m := c.ctx.Value("content")
	if res, ok := m.(JobContextMap).m["result"]; ok {
		return res
	}
	return nil
}

func NewResource(resName string, val interface{}) *Resource {
	return &Resource{
		resourceName: resName,
		res:          val,
	}
}
