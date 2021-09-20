package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var (
	DefaultDatabaseFilePath = "./database.csv"
	InfoLogger              *log.Logger
	ErrorLogger             *log.Logger
	db                      = make(database)
	exit                    = make(chan bool)
	finishedReadingInputs   = make(chan bool)
	interruptChannel        = make(chan os.Signal, 1)
)

type database map[string]*Job

// deque jobs and run them
func runJobs(finishedReadingInputs chan bool, s *Scheduler) {
	for {
		for s.Len() > 0 {
			if s.CheckIfNextJobRunnable() {
				job := s.Pop().(*Job)
				InfoLogger.Println("running job", job.name)
				// execute command in bash with arguments and flags
				process := exec.Command("bash", "-c", job.process)
				_, err := process.Output()
				if err != nil {
					InfoLogger.Printf("job %v completed\n", job.name)
				}

				InfoLogger.Println("finished running job", job.name)
			}
		}
		if <-finishedReadingInputs {
			exit <- true
		}
	}

}
func persistRemainingCommands(scheduler *Scheduler) {
	// if program was terminated, persist commands to database (file)
	<-interruptChannel
	InfoLogger.Println("Interrupt received, persisting remaining commands to database")
	dbFile, err := os.Create("unfinishedJobs.csv")
	defer dbFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(*scheduler); i++ {
		job := (*scheduler)[i]
		name := job.name
		interval := strconv.Itoa(int(job.interval.Seconds()))
		process := job.process
		_, err := fmt.Fprintf(dbFile, "%v,%v,%v\n", name, interval, process)
		if err != nil {
			log.Fatal(err)
		}
	}
	os.Exit(0)
}
func main() {
	// Setup loggers
	InfoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Setup interrupt signal handler
	signal.Notify(interruptChannel, os.Interrupt, os.Kill)

	// Open the mock database file or create one if file does not exist
	_, err := os.Stat(DefaultDatabaseFilePath)
	file, err := os.Open(DefaultDatabaseFilePath)
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	// Initialize the scheduler data structure (heap)
	schedulerObj := Scheduler{}
	scheduler := InitScheduler(schedulerObj)
	go persistRemainingCommands(scheduler)

	// load in database jobs if they exist
	fmt.Println(file.Name())
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ",")
		//name, interval, process
		name := tokens[0]
		intervalSeconds, err := strconv.Atoi(tokens[1])
		if err != nil {
			ErrorLogger.Fatal(err)
		}
		interval := time.Second * time.Duration(intervalSeconds)
		processString := tokens[2]
		job := CreateJob(name, interval, processString)
		db[name] = &job
		InfoLogger.Printf("Added job %v to database\n", name)
		scheduler.AddJob(&job)
	}
	// Start running jobs in another thread
	go runJobs(finishedReadingInputs, scheduler)

	// simultaneously add jobs from console input and schedule jobs
	fmt.Println("Reading inputs from console (press q to quit)")

	// read input from console (first line is command, second line is command metadata)
	fmt.Println("Enter your command (eg. \"go run main_test.go\")")
	inputScanner := bufio.NewScanner(os.Stdin)
	for inputScanner.Scan() {
		command := inputScanner.Text()
		if command != "q" {
			fmt.Println("Enter command metadata in the format: NAME, INTERVAL_SECONDS")
			inputScanner.Scan()
			line := strings.Split(inputScanner.Text(), " ")
			jobName := line[0]
			jobIntervalSeconds, err := strconv.Atoi(line[1])
			jobIntervalDuration := time.Second * time.Duration(jobIntervalSeconds)
			if err != nil {
				ErrorLogger.Fatal(err)
			}
			// Create job from input data and add job to scheduler
			job := CreateJob(jobName, jobIntervalDuration, command)
			scheduler.AddJob(&job)
		} else {
			// break if "q" command (quit)
			break
		}
		fmt.Println("Enter your command (eg. \"go run main_test.go\")")
	}

	// Send notification to scheduler thread that inputs are finished being read
	finishedReadingInputs <- true

	// Wait for notification from scheduler thread that all jobs have been scheduled
	<-exit
}
