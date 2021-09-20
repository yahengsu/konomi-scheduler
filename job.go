package main

import (
	"errors"
	"time"
)

type Job struct {
	name string
	interval time.Duration
	process string
	queuedTimestamp time.Time
	index int
	// timeoutTimestamp time.Time
}

func CreateJob(name string, interval time.Duration, process string) Job {
	return Job{
		name:     name,
		interval: interval,
		process:  process,
		queuedTimestamp: time.Time{},
	}
}

func DeleteJob(name string, db map[string]*Job, s *Scheduler) error {
	job, exists := db[name]
	if !exists {
		return errors.New("job does not exist in database")
	}
	if job.index != -1 {
		s.RemoveJob(job)
	}
	delete(db, name)
	return nil
}