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

	// correct duration based on MsgNumber
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

// TODO
func (b *stanBench) initPubRunner() {
	b.pubers = make([]Runner, 0, b.conf.PubNum)
	for k := 0; k < b.conf.PubNum; k++ {
		for h := 0; h < b.conf.SubjectNum; h++ {
			b.pubers = append(b.pubers, b.newPubRunner(h))
		}
	}

	return
}

// TODO
func (b *stanBench) newPubRunner(subjectIndex int) Runner {
	conn := b.createConnect()

	return NewRunner(b.ctx, RunnerOpts{
		Task:  b.newPubTask(conn, subjectIndex),
		Total: int(b.conf.MsgNumber) / int(b.conf.PubNum),
		Rate:  b.conf.Rate,
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
		Task:  b.newSubTask(conn, subjectIndex),
		Total: 1,
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

		b.reporter.Report(err,
			subject,
			len(b.data),
			time.Since(start))
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
		start := time.Now()
		select {
		case <-b.ctx.Done():
			return
		case m := <-msgCh:
			b.reporter.Report(nil,
				subject,
				len(m.Data),
				time.Since(start))
		}
	}
}

func (b *stanBench) run() error {
	b.start = time.Now()
	fmt.Printf("start running with:%v, at:%v\n", b.conf, b.start)

	b.invokeTaskRunner()
	b.reporter.Start()

	fmt.Println("waiting for complete")

	<-time.After(b.conf.Duration)

	b.cancel()
	b.end = time.Now()
	fmt.Printf("running complete at:%v\n\n\tElapsed:%v\n", time.Now(), b.end.Sub(b.start))

	b.wg.Wait()
	// FIXME: reporter not work
	fmt.Printf("Report:%v\n", b.reporter.GetReport())

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

	log.Printf("runner number: %v", count)
}

func (b *stanBench) startGoroutine(f func()) {
	b.wg.Add(1)
	go func() {
		f()
		log.Printf("runner done")
		b.wg.Done()
	}()
}
