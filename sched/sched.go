package sched

import (
	"context"
	"errors"
)

func NewCaptin() *Captin {
	c := Captin{
		schedulers: make(map[string]*Scheduler),
	}
	return &c
}

func (s *Scheduler) worker(jobs chan Job, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.wg.Done()
			return
		case j := <-jobs:
			if j.args == nil {
				j.args = struct{}{}
			}
			j.ref.handler(j.earlyBreak, j.resChan, j.args)
			if j.ref.hasRV {
				close(j.resChan)
			}
			if j.ref.preemptable {
				close(j.earlyBreak)
			}
		default:
		}
	}
}

func (c *Captin) NewScheduler(schedName string, workerNum int) (*Scheduler, error) {
	newS := &Scheduler{
		schedName: schedName,
		workerNum: workerNum,
		jobMap:    make(map[int][]*Job),
		jobDefs:   make(map[string]*JobDef),
		prioMap:   make(map[int]int),
		recvChan:  make(chan Job, 100),
		jobQueue:  make(chan Job, 100),
	}
	if oldS, ok := c.schedulers[schedName]; ok {
		return oldS, errors.New("SchedName exists under this captin")
	} else {
		c.schedulers[schedName] = newS
	}

	return newS, nil
}

func (c *Captin) FindScheduler(schedName string) (*Scheduler, error) {
	if s, ok := c.schedulers[schedName]; ok {
		return s, nil
	}
	return nil, errors.New("Scheduler not found")
}

func (s *Scheduler) Start(ctx context.Context) {
	for i := 0; i < priorityGap; i++ {
		s.prioMap[i] = 0
	}

	s.wg.Add(s.workerNum)
	for i := 0; i < s.workerNum; i++ {
		go s.worker(s.jobQueue, ctx)
	}
	for {
		select {
		case <-ctx.Done():
			close(s.jobQueue)
			close(s.recvChan)
			s.wg.Wait()
			return
		case j := <-s.recvChan:
			s.prioMap[j.ceilPriority]++
			jobList := s.jobMap[j.ceilPriority]
			s.jobMap[j.ceilPriority] = append(jobList, &j)
			// TODO: ceiling the prioriy of specific job,
			// if it has the sync resource shared with this job
		default:
			// TODO: watchDog
			s.sched()
		}
	}
}

func (s *Scheduler) sched() {
	var hPrio int
	for hPrio = 0; hPrio < priorityGap; hPrio++ {
		if s.prioMap[hPrio] > 0 {
			break
		}
	}
	if len(s.jobMap[hPrio]) > 0 {
		var j *Job
		jobList := s.jobMap[hPrio]
		j = jobList[0]
		s.jobQueue <- *j
		if len(jobList) > 1 {
			s.jobMap[hPrio] = jobList[1:]
		} else {
			s.jobMap[hPrio] = make([]*Job, 0)
		}
		s.prioMap[j.ref.priority]--
	}
}
