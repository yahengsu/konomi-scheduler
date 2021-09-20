# Konomi Network - Design a Job Scheduler

A job scheduler developed in Go. Supports the following:

- Create/Read/Update/Delete jobs to run
- Support one time execution and repetitive executions at a fixed interval
- Jobs are persisted

## Details

- Persisted jobs are loaded from our database and added to the scheduler
- Scheduler then queues/dequeues jobs on a seperate thread and our main thread waits for jobs to be entered from the CLI
  - Using CLI for adding jobs is not ideal - preferably would have liked some HTTP API implementation or
  - Scheduler data structure is a Priority Queue (Heap) where jobs that are scheduled to run earlier are prioritized
- Jobs are run on their own thread using the exec library
- On program shutdown, persists remaining jobs in scheduler queue to database (file)

## Running the Program

- Clone the repository
- Run `docker build --tag docker-konomi-project .`
- Run `docker run docker-konomi-project`

Or:

- Clone the repository
- Run `go build`
- Run the executable `./konomi-project` or `./konomi-project.exe`

### Create/Read/Update/Delete Jobs

- Uses a csv file to mock database functionality
  - Tabular data (mock SQL database), easier to represent data as csv for mock
- Creating a job creates a job struct which contains metadata about the job which are stored in a Map
  - To create/read/update jobs we just need to modify our Job in the Map
- Jobs are persisted through the text file (mock database)

### Support One Time and Repetitive Executions

- User can set the execution type and repetition interval during job creation, can also update these settings by interacting with the map container directly
- For repetitive executions, add a new job to the scheduler with the specified runtime timestamp such that the job will only run after the current time has passed that timestamp

## Other Considerations

### Timeout Detection

Timeout detection can be taken care of by adding a timeout field to the Job struct, containing the timestamp
of when the Job is expected to finish. Since the Job will be removed by the scheduler once it is finished, we will know that we have to remove the job if the current timestamp is greater than the timeout timestamp. A naive implementation would be to scan the job queue every second or so and check if there are jobs that have timed out and need to be removed. A better solution would be to set up a notification system that would only check the job queue at the timeout timestamps and see if they need to be removed instead of every second, reducing the amount of processing power wasted.

### Retries

If a job reaches its timeout timestamp or if an error code is returned from the process then we can reschedule the job with the exact same parameters and add it to the scheduler queue again.

### Improving Database Performance

We can maybe add a cache layer in RAM by batching our database calls and fetching the next N jobs that need to be run and storing them to our scheduler structure instead of making a database call every time we need to fetch a job to run.

### Scalability

We can scale up this solution to thousands of jobs and many workers by creating a different ThreadPool for each main job function (scheduler, input handler, process executor) and adding more threads for each system when needed. Some synchronization will be required to handle communication between different threads, but the scheduling and input handling should be very paralellizable.
