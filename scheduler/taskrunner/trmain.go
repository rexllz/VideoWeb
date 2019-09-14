package taskrunner

import "time"

type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

func NewWorker(interval time.Duration, r *Runner) *Worker {

	return &Worker{
		ticker: time.NewTicker(interval * time.Second),
		runner: r,
	}
}

func (w *Worker) startWorker() {

	//not good , this (range) is a synchronization method
	//for c = range w.ticker.C {
	//
	//}

	println("start worker!")
	for  {
		select {
		// get ticker 's channel signal
		case <- w.ticker.C:
			go w.runner.StartAll()
		}
	}
}

func Start()  {

	println("start!")
	r := NewRunner(3, true, VideoClearDispatcher, VideoClearExecutor)
	w := NewWorker(3, r)
	go w.startWorker()

}