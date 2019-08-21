package report

import (
	"context"
	"fmt"
	"time"

	"github.com/shohi/nats-bench/pkg/metrics"
)

// Status represents Reporter's state
type Status int

const (
	Ready Status = iota
	Running
	Stopped
)

// Reporter ...
type Reporter struct {
	ctx    context.Context
	opCh   chan func() // operation chan
	status Status

	dumpInterval time.Duration
	stats        Stats
}

// NewReporter ...
func NewReporter(ctx context.Context, reportInterval time.Duration) *Reporter {
	reporter := &Reporter{
		ctx:    ctx,
		opCh:   make(chan func(), 128),
		status: Ready,

		dumpInterval: reportInterval,
	}

	return reporter
}

// Start ...
func (r *Reporter) Start() {
	r.stats.Start(time.Now())

	go func() {
		ticker := time.NewTicker(r.dumpInterval)
		for {
			select {
			case <-r.ctx.Done():
				r.status = Stopped
				ticker.Stop()
			case op := <-r.opCh:
				if op == nil {
					r.status = Stopped
					ticker.Stop()
					return
				}

				op()
			case <-ticker.C:
				fmt.Println(r.stats.Snapshot())
			}
		}
	}()
}

// NOTE: not goroutine safe, should only be accessed in loop goroutine.
func (r *Reporter) getStatus() Status {
	return r.status
}

// GetReport from reporter.
func (r *Reporter) GetReport() string {
	ch := make(chan string, 1)
	r.opCh <- func() {
		ch <- r.stats.Final()
	}

	return <-ch
}

// Report a sample
func (r *Reporter) Report(t metrics.Txn) {
	r.opCh <- func() {
		if r.getStatus() == Stopped {
			return
		}

		r.stats.AddTxn(t)
	}
}

func (r *Reporter) Stop() {
	r.opCh <- nil
}
