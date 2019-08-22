package bench

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"

	"github.com/shohi/nats-bench/pkg/config"
	"github.com/shohi/nats-bench/pkg/metrics"
	"github.com/shohi/nats-bench/pkg/report"
)

// RunStanBench start benchmark test for nats-streaming server
func RunStanBench(conf config.StanConfig) error {
	b := newStanBench(conf)
	return b.run()
}

type stanBench struct {
	conf config.StanConfig

	ctx    context.Context
	cancel context.CancelFunc

	wg     sync.WaitGroup
	pubers []Runner
	subers []Runner

	data     []byte
	reporter *report.Reporter

	start time.Time
	end   time.Time
}

func newStanBench(conf config.StanConfig) *stanBench {
	fakeData := make([]byte, conf.MsgSize)

	var ctx context.Context
	var cancel context.CancelFunc

	// TODO: correct duration based on MsgNumber
	if conf.Duration == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), conf.Duration)
	}

	reporter := report.NewReporter(ctx, conf.ReportInterval)

	b := &stanBench{
		conf: conf,

		data:     fakeData,
		ctx:      ctx,
		cancel:   cancel,
		reporter: reporter,
	}

	if b.conf.Duration > 0 {
		b.conf.MsgNumber = 0
	}

	b.initRunner()

	return b
}

func (b *stanBench) createConnect() stan.Conn {
	nc, err := nats.Connect(b.conf.URL, nats.Timeout(b.conf.NatsConnectTimeout))
	if err != nil {
		panic(err)
	}

	conn, err := stan.Connect(b.conf.Cluster,
		NewID("stan-client"),
		stan.ConnectWait(b.conf.StanConnectTimeout),
		stan.PubAckWait(b.conf.PubAckWait),
		stan.MaxPubAcksInflight(b.conf.PubAckMaxInflight),
		stan.NatsConn(nc),
		stan.Pings(b.conf.PingInterval, b.conf.PingMax),
	)
	if err != nil {
		panic(err)
	}

	return conn
}

func (b *stanBench) initRunner() {
	b.initPubRunner()
	b.initSubRunner()
}

func (b *stanBench) initPubRunner() {
	b.pubers = make([]Runner, 0, b.conf.PubNum)
	for k := 0; k < b.conf.PubNum; k++ {
		for h := 0; h < b.conf.SubjectNum; h++ {
			b.pubers = append(b.pubers, b.newPubRunner(h))
		}
	}

	return
}

func (b *stanBench) newPubRunner(subjectIndex int) Runner {
	conn := b.createConnect()

	return NewRunner(b.ctx, RunnerOpts{
		Task:   b.newPubTask(conn, subjectIndex),
		Total:  int(b.conf.MsgNumber) / int(b.conf.PubNum),
		Rate:   b.conf.Rate,
		OnDone: func() { _ = conn.Close() },
	})

}

func (b *stanBench) initSubRunner() {
	b.subers = make([]Runner, 0, b.conf.SubNum)

	for k := 0; k < b.conf.SubNum; k++ {
		for h := 0; h < b.conf.SubjectNum; h++ {
			b.subers = append(b.subers, b.newSubRunner(h))
		}
	}

	return
}

func (b *stanBench) newSubRunner(subjectIndex int) Runner {
	conn := b.createConnect()

	return NewRunner(b.ctx, RunnerOpts{
		Task:   b.newSubTask(conn, subjectIndex),
		Total:  1,
		OnDone: func() { _ = conn.Close() },
	})

}

func (b *stanBench) newPubTask(conn stan.Conn, subjectIndex int) Task {
	subject := fmt.Sprintf("%s-%v", b.conf.SubjectPrefix, subjectIndex)

	return func() {
		start := time.Now()
		var err error
		if b.conf.PubAsync {
			_, err = conn.PublishAsync(subject, b.data, nil)
		} else {
			err = conn.Publish(subject, b.data)
		}

		b.reporter.Report(metrics.Txn{
			Name:  subject,
			Err:   err,
			Size:  float64(len(b.data)),
			Start: start,
			End:   time.Now(),
		})
	}
}

func (b *stanBench) newSubTask(conn stan.Conn, subjectIndex int) Task {
	subject := fmt.Sprintf("%s-%v", b.conf.SubjectPrefix, subjectIndex)

	msgCh := make(chan *stan.Msg, 1024)
	_, err := conn.QueueSubscribe(
		subject,
		b.conf.QueueGroup, func(m *stan.Msg) {
			msgCh <- m
		},
		stan.StartWithLastReceived())

	if err != nil {
		panic(err)
	}

	return func() {
		// start := time.Now()
		select {
		case <-b.ctx.Done():
			msgCh <- nil
		case m := <-msgCh:
			// if context is done, exit task.
			if m == nil {
				return
			}

			/*
				// NOTE: not report subscriber's stat.
				b.reporter.Report(metrics.Txn{
					Name:  subject,
					Err:   nil,
					Size:  float64(len(m.Data)),
					Start: start,
					End:   time.Now(),
				})
			*/
		}
	}
}

func (b *stanBench) run() error {
	b.start = time.Now()
	log.Printf("start running with: [%+v]\n", b.conf)

	b.invokeTaskRunner()
	b.reporter.Start()

	<-time.After(b.conf.Duration)

	b.cancel()
	b.end = time.Now()
	log.Printf("completed. Elapsed: %v\n", b.end.Sub(b.start))

	b.wg.Wait()
	log.Printf("Report:\n\n%v\n", b.reporter.GetReport())

	return nil
}

func (b *stanBench) invokeTaskRunner() {
	count := 0
	// Publishers
	for _, r := range b.pubers {
		count++
		b.startGoroutine(r.Run)
	}

	// Subscribers
	for _, r := range b.subers {
		count++
		b.startGoroutine(r.Run)
	}
}

func (b *stanBench) startGoroutine(f func()) {
	b.wg.Add(1)
	go func() {
		f()
		b.wg.Done()
	}()
}
