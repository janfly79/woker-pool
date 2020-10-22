package logic

import (
	"fmt"
	"time"
)

var (
	MaxWorker = 10
	MaxQueue  = 100
)

type PayloadCollection struct {
	WindowsVersion string    `json:"version"`
	Token          string    `json:"token"`
	Payloads       []Payload `json:"data"`
}

type Payload struct {
	// [redacted]
	Now time.Time
}

func (p *Payload) UploadToS3() error {
	// the storageFolder method ensures that there are no name collision in
	// case we get same timestamp in the key name
	storage_path := fmt.Sprintf("%v/%v", "accd", time.Now().UnixNano())

	fmt.Println("上传成功到 ", storage_path)

	time.Sleep( 3 * time.Second)

	return nil
}

// Job represents the job to be run
type Job struct {
	Payload Payload
}

// A buffered channel that we can send work requests on.
var JobQueue chan Job = make(chan Job, 20)



func NewJob(){
	fmt.Println("new job push start")
	payload := Payload{Now:time.Now()}
	work := Job{Payload: payload}

	// Push the work onto the queue.
	JobQueue <- work

	fmt.Println("new job push end")

}


// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				if err := job.Payload.UploadToS3(); err != nil {
					fmt.Printf("Error uploading to S3: %s", err.Error())
				}
				fmt.Printf("job done now : %v ms \n", time.Since(job.Payload.Now).Milliseconds())

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
