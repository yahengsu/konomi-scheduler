package main

import (
	"errors"
	"time"
)

// Job - Represents a job
type Job struct {
	name            string
	interval        time.Duration
	process         string
	queuedTimestamp time.Time
	index           int
	// timeoutTimestamp time.Time
}

// CreateJob - creates a new Job struct and returns it
func CreateJob(name string, interval time.Duration, process string) Job {
	return Job{
		name:            name,
		interval:        interval,
		process:         process,
		queuedTimestamp: time.Time{},
	}
}

// DeleteJob - Deletes a job from the database and the scheduler if it is queued
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
