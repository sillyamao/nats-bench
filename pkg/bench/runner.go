package bench

import (
	"context"
	"log"
	"time"
)

// TODO: add context support for task
type Task func()

// TODO: use golang/x/time/rate
type RunnerOpts struct {
	Task  Task
	Total int
	Rate  float64 // number of Task should be performed in one second. If zero, means no limit.
}

type runner struct {
	ctx  context.Context
	opts RunnerOpts

	done      int
	interval  time.Duration
	isRunning bool
}

func NewRunner(ctx context.Context, opts RunnerOpts) Runner {
	var interval time.Duration
	// no limit is set for zero or negative value
	if opts.Rate <= 0 {
		interval = 0
	} else {
		interval = time.Duration(float64(time.Second) / opts.Rate)
	}
	if opts.Total <= 0 {
		opts.Total = 0
	}

	return &runner{
		ctx:      ctx,
		opts:     opts,
		interval: interval,
	}
}

func (r *runner) Run() {
	if r.isRunning {
		panic("already running")
	}
	r.isRunning = true

	for r.shouldRun() {
		select {
		case <-r.ctx.Done():
			log.Printf("runner exit")
			return
		default:
			used := elapsed(r.opts.Task)
			if used < r.interval {
				time.Sleep(r.interval - used)
			}
		}
	}
}

func (r *runner) shouldRun() bool {
	if r.opts.Total == 0 {
		return true
	}

	if r.done < r.opts.Total {
		r.done++
		return true
	}

	return false
}

func elapsed(run func()) time.Duration {
	s := time.Now()
	run()
	return time.Since(s)
}
