package main

import (
	"container/heap"
	"time"
)

// Scheduler - Data structure for queueing/de-queueing jobs
type Scheduler []*Job

func (s Scheduler) Len() int {
	return len(s)
}

func (s Scheduler) Less(i, j int) bool {
	return s[i].queuedTimestamp.Unix() > s[j].queuedTimestamp.Unix()
}

func (s Scheduler) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *Scheduler) Push(x interface{}) {
	n := len(*s)
	job := x.(*Job)
	job.index = n
	job.queuedTimestamp = time.Now()
	*s = append(*s, job)
}

func (s *Scheduler) Pop() interface{} {
	old := *s
	n := len(old)
	job := old[n-1]
	old[n-1] = nil
	job.index = -1
	*s = old[0 : n-1]
	for _, otherJobs := range *s {
		otherJobs.index -= 1
	}
	// if recurring task then we need to add it to scheduler again so it stays recurring
	if job.interval != (time.Second * time.Duration(0)) {
		recurringJob := Job{name: job.name, interval: job.interval, process: job.process}
		recurringJob.queuedTimestamp = time.Now()
		heap.Push(s, &recurringJob)
	}
	return job
}

func (s *Scheduler) AddJob(job *Job) {
	heap.Push(s, job)
}

func (s *Scheduler) RemoveJob(job *Job) {
	heap.Remove(s, job.index)

}

func (s *Scheduler) CheckIfNextJobRunnable() bool {
	// case where we have only jobs at recurring intervals, want to schedule after interval_time has passed
	n := s.Len()
	job := (*s)[n-1]
	if job.queuedTimestamp.Add(job.interval).Unix() > time.Now().Unix() {
		return false
	}
	return true
}

func InitScheduler(s Scheduler) *Scheduler {
	h := &s
	heap.Init(h)
	return h
}
