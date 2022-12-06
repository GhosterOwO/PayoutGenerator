package dispatch

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

var MaxJobs = int64(0)

type Dispatcher struct {
	WorkerPool chan chan Job
	JobQueue   chan Job
	Complete   chan Job
	JobTotal   int64
	JobDone    int64
	Jobs       int64
	wg         sync.WaitGroup
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		JobQueue:   make(chan Job),
		Complete:   make(chan Job),
		WorkerPool: make(chan chan Job, maxWorkers),
		wg:         sync.WaitGroup{},
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < cap(d.WorkerPool); i++ {
		worker := NewWorker(d.WorkerPool, d.Complete)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) block() {
	for d.Jobs >= MaxJobs && MaxJobs > 0 {
		select {
		case <-time.After(time.Millisecond):
			continue
		}
	}
}

func (d *Dispatcher) Join(work Job) {
	d.block()
	d.wg.Add(1)
	atomic.AddInt64(&d.Jobs, 1)
	go func() {
		atomic.AddInt64(&d.JobTotal, 1)
		d.JobQueue <- work
	}()
}

func (d *Dispatcher) Progress() float64 {
	if d.JobDone == 0 {
		if d.JobTotal == 0 {
			return float64(100.0)
		}
		return float64(0.0)
	}
	ratio := float64(d.JobDone) / float64(d.JobTotal)
	return math.Min(ratio*float64(100.0), float64(100.0))
}

func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		case <-d.Complete:
			atomic.AddInt64(&d.JobDone, 1)
			atomic.AddInt64(&d.Jobs, -1)
			d.wg.Done()
		}
	}
}
