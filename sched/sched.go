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
			j.ref.handler(j.earlyBreak, j.args)

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
	index := 0
	for {
		select {
		case <-ctx.Done():
			close(s.jobQueue)
			close(s.recvChan)
			s.wg.Wait()
			return
		case j := <-s.recvChan:
			s.prioMap[j.ref.priority]++
			jobList := s.jobMap[j.ref.ceilPriority]
			s.jobMap[j.ref.ceilPriority] = append(jobList, &j)
			index = min(j.ref.priority, index)
			// TODO: ceiling the prioriy of specific job,
			// if it has the sync resource shared with this job
		default:
			// TODO: watchDog
			s.sched(index)
			index++
			// reset index
			if index > priorityGap {
				index = 0
			}
		}
	}
}

func (s *Scheduler) sched(index int) {
	if s.prioMap[index] > 0 && len(s.jobMap[index]) > 0 {
		var j *Job
		jobList := s.jobMap[index]
		j = jobList[0]
		s.jobQueue <- *j
		if len(jobList) > 1 {
			s.jobMap[index] = jobList[1:]
		} else {
			s.jobMap[index] = make([]*Job, 0)
		}
		s.prioMap[j.ref.priority]--
		if j.ref.preemptable {
			close(j.earlyBreak)
		}
	}
}
