package dispatch

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	complete   chan Job
	quit       chan bool
}

func NewWorker(pool chan chan Job, complete chan Job) Worker {
	return Worker{
		WorkerPool: pool,
		complete:   complete,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				job.Execute()
				w.complete <- job
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
