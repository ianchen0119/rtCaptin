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

func (s *Scheduler) worker(jobs chan *Job, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.wg.Done()
			return
		case j := <-jobs:
			ctx := j.NewJobContext()
			j.ref.handler(ctx, j.args)
			if j.ref.hasRV && j.resChan != nil {
				close(j.resChan)
			}
			if j.ref.preemptable && j.earlyBreak != nil {
				close(j.earlyBreak)
			}
			j.done = true
		default:
		}
	}
}

func (c *Captin) NewScheduler(schedName string, workerNum int) (*Scheduler, error) {
	newS := &Scheduler{
		schedName: schedName,
		workerNum: workerNum,
		jobMap:    []*Job{},
		jobDefs:   make(map[string]*JobDef),
		recvChan:  make(chan Job, workerNum*10),
		jobQueue:  make(chan *Job, workerNum*10),
		ceilChan:  make(chan *Job, 1),
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
	s.wg.Add(s.workerNum)
	for i := 0; i < s.workerNum; i++ {
		go s.worker(s.jobQueue, ctx)
	}
	for {
		// The event, sent from ohther scheduler, has highest priority.
		// The Scheduler needs to proceed this request first.
		select {
		case j := <-s.ceilChan:
			j.owner.ceilPriority(j)
		default:
		}
		select {
		case <-ctx.Done():
			close(s.jobQueue)
			close(s.recvChan)
			s.wg.Wait()
			return
		case j := <-s.recvChan:
			s.jobMap = append(s.jobMap, &j)
			if j.resources != nil {
				for i := 0; i < len(j.resources); i++ {
					resource := j.resources[i]
					if resource.belong == nil {
						resource.belong = &j
						continue
					} else {
						if resource.belong.owner != s {
							resource.belong.owner.ceilPriorityToOhterScheder(resource.belong)
							resource.belong = &j
							continue
						}
					}
					curJ := resource.belong
					if curJ != nil || !curJ.done {
						if curJ.ceilPriority > j.ref.priority {
							resource.belong = &j
							// incomming job is more important than the job
							// that uses the sync resource
							if curJ.ref.preemptable {
								// if current job is preemptable,
								// terminate this job by sending the signal
								curJ.earlyBreak <- 1
								continue
							}
							if curJ != nil && curJ != &j {
								// otherwise, raise the priority of current job
								s.ceilPriority(curJ)
							}
						}
					}
				}
			}
		default:
			// TODO: watchDog
			s.sched()
		}
	}
}

func (s *Scheduler) ceilPriority(cJob *Job) {
	cJob.ceilPriority = 0
}

func (s *Scheduler) ceilPriorityToOhterScheder(j *Job) {
	j.owner.ceilChan <- j
}

// Earliest-Deadline-First (EDF) Algorithm
func (s *Scheduler) sched() {
	var selectedJobIndex int = 0
	var prioMin = 256
	for i := 0; i < len(s.jobMap); i++ {
		if ok, selectedNum := min(s.jobMap[i].ceilPriority, prioMin); ok {
			prioMin = selectedNum
			selectedJobIndex = i
		}
		// aging
		if s.jobMap[i].ceilPriority > 0 {
			s.jobMap[i].ceilPriority--
		}
	}
	if len(s.jobMap) > 0 {
		j := s.jobMap[selectedJobIndex]
		s.jobQueue <- j
		s.jobMap = remove(s.jobMap, selectedJobIndex)
	}
}
