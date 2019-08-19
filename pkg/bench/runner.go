package bench

import "time"

type Task func()

type RunnerOpts struct {
	Task      Task
	Total     int
	Frequency float64 // number of requests sent in one second
}

type runner struct {
	opts RunnerOpts

	done      int
	interval  time.Duration
	isRunning bool
}

func NewRunner(opts RunnerOpts) Runner {
	interval := time.Duration(float64(time.Second) / opts.Frequency)

	return &runner{
		opts:     opts,
		interval: interval,
	}
}

func (r *runner) Stop() {
	if !r.isRunning {
		panic("not running")
	}

	// TODO

}

func (r *runner) Run() {
	if r.isRunning {
		panic("already running")
	}
	r.isRunning = true

	for r.done = 0; r.done < r.opts.Total; r.done++ {
		used := elapsed(r.opts.Task)
		if used < r.interval {
			time.Sleep(r.interval - used)
		}
	}
}

func elapsed(run func()) time.Duration {
	s := time.Now()
	run()
	return time.Since(s)
}
