package main

import (
	"fmt"
	"testing"
	"time"
)

// TestCreateJob - Tests the CreateJob method
func TestCreateJob(t *testing.T) {
	name := "name"
	interval := time.Second * time.Duration(0)
	process := "echo hello"
	job := CreateJob(name, interval, process)

	if job.name != name || job.interval != interval || job.process != process {
		t.Errorf("Want name=%v, interval=%v, process=%v, got name=%v, interval=%v, process=%v", name, interval, process, job.name, job.interval, job.process)
	}
}

// TestReadJob - Tests the ReadJob method
func TestReadJob(t *testing.T) {
	db := make(database)
	name := "name"
	interval := time.Second * time.Duration(0)
	process := "echo hello"
	job := CreateJob(name, interval, process)
	db[name] = &job
	val, ok := db[name]
	if !ok {
		t.Errorf("Job was not found in map")
	}

	if val.name != job.name || val.interval != job.interval || val.process != job.process {
		t.Errorf("Jobs do not match: Want name=%v, interval=%v, process=%v, got name=%v, interval=%v, process=%v", name, interval, process, val.name, val.interval, val.process)
	}

}

// TestUpdateJob - Tests the UpdateJob method
func TestUpdateJob(t *testing.T) {
	db := make(database)
	name := "name"
	interval := time.Second * time.Duration(0)
	process := "echo hello"
	newProcess := "echo hi"
	job := CreateJob(name, interval, process)
	db[name] = &job
	db[name].process = newProcess
	if job.process != newProcess {
		t.Errorf("Expected job.process=%v, got job.process=%v", newProcess, job.process)
	}
}

// TestDeleteJob - Tests the DeleteJob method
func TestDeleteJob(t *testing.T) {
	db := make(database)
	s := Scheduler{}
	sch := InitScheduler(s)
	name := "name"
	interval := time.Second * time.Duration(0)
	process := "echo hello"
	job := CreateJob(name, interval, process)
	db[name] = &job
	sch.AddJob(&job)
	err := DeleteJob(name, db, sch)
	if err != nil {
		t.Errorf("Error occurred while deleting job: %v index=%v", err, job.index)
	}
}

// TestScheduler_AddJob - Tests the AddJob method of the Scheduler
func TestScheduler_AddJob(t *testing.T) {
	s := Scheduler{}
	sch := InitScheduler(s)
	name := "name"
	interval := time.Second * time.Duration(0)
	process := "echo hello"
	job := CreateJob(name, interval, process)
	sch.AddJob(&job)
	addedJob := sch.Pop().(*Job)

	if job.name != addedJob.name || job.interval != addedJob.interval ||
		job.process != addedJob.process || job.index != addedJob.index ||
		job.queuedTimestamp != addedJob.queuedTimestamp {
		t.Errorf("Jobs do not match: Want name=%v, interval=%v, process=%v, index=%v, queuedTimestamp=%v, "+
			"got name=%v, interval=%v, process=%v, index=%v, queuedTimestamp=%v",
			job.name, job.interval, job.process, job.index, job.queuedTimestamp,
			addedJob.name, addedJob.interval, addedJob.process, addedJob.index, addedJob.queuedTimestamp)
	}
}

// TestSceduler_RemoveJob - Tests the RemoveJob method of the Scheduler
func TestScheduler_RemoveJob(t *testing.T) {
	s := Scheduler{}
	sch := InitScheduler(s)
	name := "name"
	interval := time.Second * time.Duration(0)
	process := "echo hello"
	job := CreateJob(name, interval, process)

	sch.AddJob(&job)
	fmt.Println(sch.Len())

	fmt.Println(sch.Len(), job.index)
	sch.RemoveJob(&job)
	if sch.Len() != 0 {
		t.Errorf("Expected scheduler len=0, got len=%v", sch.Len())
	}
}

// TestScheduler_CheckIfNextJobRunnable - Tests the CheckIfNextJobRunnable method of the Scheduler
func TestScheduler_CheckIfNextJobRunnable(t *testing.T) {
	s := Scheduler{}
	sch := InitScheduler(s)
	name := "name"
	interval := time.Second * time.Duration(0)
	process := "echo hello"
	job := CreateJob(name, interval, process)
	sch.AddJob(&job)
	if !sch.CheckIfNextJobRunnable() {
		t.Errorf("Expected true got false")
	}
	job.queuedTimestamp = job.queuedTimestamp.Add(time.Minute)

	if sch.CheckIfNextJobRunnable() {
		t.Errorf("Expected false got true, timestamp=%v", job.queuedTimestamp)
	}
}
