package bench

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/stan.go"

	"github.com/shohi/nats-bench/pkg/config"
	"github.com/shohi/nats-bench/pkg/report"
)

// RunStanBench start benchmark test for nats-streaming server
func RunStanBench(conf config.StanConfig) error {
	b := newStanBench(conf)
	return b.run()
}

// TODO
type stanBench struct {
	conf config.StanConfig

	ctx      context.Context
	cancelFn context.CancelFunc

	runners  []Runner
	data     []byte
	reporter *report.Reporter
}

// TODO
func newStanBench(conf config.StanConfig) *stanBench {
	fakeData := make([]byte, conf.MsgSize)

	ctx, cancel := context.WithCancel(context.Background())
	reporter := report.NewReporter(ctx, conf.ReportInterval)

	return &stanBench{
		conf: conf,

		data:     fakeData,
		ctx:      ctx,
		cancelFn: cancel,
		reporter: reporter,
	}
}

func (b *stanBench) createConnect() (stan.Conn, error) {
	// TODO
	return nil, nil
}

func (b *stanBench) initRunner() {
	// TODO:
	// 1. init pub runner
	// 2. init sub runner
}

// TODO
func (b *stanBench) run() error {
	start := time.Now()
	fmt.Printf("start running with:%v, at:%v\n", b.conf, start)

	b.initRunner()
	b.reporter.Start()

	var wg sync.WaitGroup
	wg.Add(len(b.runners))
	for _, r := range b.runners {
		go func(x Runner) {
			x.Run()
			wg.Done()
		}(r)
	}

	fmt.Println("waiting for complete")
	select {
	case <-time.After(b.conf.Duration):
		for _, r := range b.runners {
			r.Stop()
		}
		fmt.Printf("running complete at:%v\n\n\tElapsed:%v\n", time.Now(), time.Since(start))
	}

	fmt.Printf("Report:%v\n", b.reporter.GetReport())
	b.cancelFn()

	return nil
}
