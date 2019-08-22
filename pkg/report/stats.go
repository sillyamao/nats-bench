package report

import (
	"fmt"
	"time"

	"github.com/shohi/nats-bench/pkg/metrics"
)

// Stats holds benchmark result.
type Stats struct {
	start  time.Time
	report Report // whole report

	snapshotCount int
	snapshotStart time.Time
	clientMap     map[string]*Sample
}

func (s *Stats) Start(t time.Time) {
	s.start = t
}

func (s *Stats) Snapshot() string {
	s.snapshotCount++

	r := s.generateSnapshotReport()
	s.report.Merge(r)

	outputFmt := "Total:    %v times, at:%v\n%v\nRecent:%v"

	return fmt.Sprintf(outputFmt, s.snapshotCount,
		time.Now(), s.report, r)
}

func (s *Stats) generateSnapshotReport() Report {
	r := Report{
		Start:       s.snapshotStart,
		Elapsed:     time.Since(s.snapshotStart),
		Concurrency: len(s.clientMap),
	}

	samples := make([]*Sample, 0, len(s.clientMap))
	for _, v := range s.clientMap {
		samples = append(samples, v)
	}

	// populate report and update whole report
	r.Populate(&SampleGroup{Samples: samples})

	// NOTE: reset clientMap for next snapshot
	s.resetClientMap()

	return r
}

// TODO: double check
func (s *Stats) resetClientMap() {
	for _, v := range s.clientMap {
		v.Reset()
	}
}

// Final returns final report.
func (s *Stats) Final() string {
	r := s.generateSnapshotReport()
	s.report.Merge(r)

	outputFmt := "Total:    %v times, at:%v\n%v\nRecent:%v\n"

	return fmt.Sprintf(outputFmt,
		s.snapshotCount, time.Now(),
		s.report, r)
}

func (s *Stats) AddTxn(t metrics.Txn) {
	sample := s.getSample(t.Name)
	sample.TxnSet = append(sample.TxnSet, t)
}

func (s *Stats) getSample(name string) *Sample {
	if s.clientMap == nil {
		s.snapshotStart = time.Now()
		s.clientMap = make(map[string]*Sample, 4096)
	}

	sample, ok := s.clientMap[name]
	if ok {
		return sample
	}

	sample = &Sample{
		TxnSet: make(metrics.TxnSet, 0, 1<<16),
	}
	s.clientMap[name] = sample

	return sample
}
