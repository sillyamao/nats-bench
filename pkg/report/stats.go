package report

import (
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
	panic("not implement yet")

	/*
		s.snapshotCount++

		r := Report{
			Start:       s.snapshotStart,
			Elapsed:     time.Since(s.snapshotStart),
			Concurrency: len(s.clientMap),
		}

		var us Samples
		for key, value := range r.Labels {
			labels = append(labels, key)
			r.Failed += value.Errors
			r.Completed += len(value.Samples)
			us.Append(value.Samples)
		}
		report.Report(us)
		r.TotalReport.Merge(report)
		return fmt.Sprintf("Total:    %v times, at:%v\n%v\nRecent:%v",
			r.Times,
			time.Now(),
			r.TotalReport,
			report,
		)
	*/
}

// Final returns final report.
func (s *Stats) Final() string {
	panic("TODO")
}

// FIMXE:
func (s *Stats) AddTxn(t metrics.Txn) {
	panic("Not complete yet")
	/*
		sample := s.getSample(t.Name)
		sample.TxnSet = append(sample.TxnSet, &t)
	*/
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

/*
// Report build on sameples
func (r *Report) Report(us Samples) string {
	sort.Sort(us)
	if us.Len() > 0 {
		r.Total = us.Total()
		r.TotalFlags = us.TotalFlags()
		r.Completed = us.Len()
		r.Min = us.Min()
		r.Max = us.Max()
		r.Avg = us.Avg()
		r.T90 = us.At(0.9)
		r.T99 = us.At(0.99)
		r.T999 = us.At(0.999)
	}
	return r.Print()
}

// Build builds report and return its content.
func (s *Stats) Build() string {
	s.snapshotCount++

	r := Report{}
	report.At = r.lastBuildTime
	report.Elapsed = time.Since(report.At)
	report.Concurrency = len(r.Labels)
	labels := make([]string, len(r.Labels))[:0]
	var us Samples
	for key, value := range r.Labels {
		labels = append(labels, key)
		report.Failed += value.Errors
		report.Completed += len(value.Samples)
		us.Append(value.Samples)
	}
	report.Report(us)
	r.TotalReport.Merge(report)
	return fmt.Sprintf("Total:    %v times, at:%v\n%v\nRecent:%v",
		r.Times,
		time.Now(),
		r.TotalReport,
		report,
	)
}

// BuildTotal build total report
func (r *Reports) BuildTotal() string {
	report := Report{}
	report.At = r.Start
	report.Elapsed = time.Since(report.At)
	report.Concurrency = len(r.Labels)
	labels := make([]string, len(r.Labels))[:0]
	var us Samples
	for key, value := range r.Labels {
		labels = append(labels, key)
		report.Failed += value.Errors
		report.Completed += len(value.Samples)
		us.Append(value.Samples)
		// for _, u := range value.Samples {
		// 	report.ElapsedTotal += u.Elapsed
		// 	report.TotalTransfer += u.Flag
		// }
	}
	// sort.Strings(labels)
	// report.Labels = labels
	report.Report(us)
	r.TotalReport.Merge(report)

	return fmt.Sprintf("Reports:%v", r.TotalReport.Print())
}

// BuildAndClear build report and clear all snapshots.
func (s *Stats) BuildAndClear() string {
	report := r.Build()
	r.Start = time.Now()
	for _, value := range r.Labels {
		value.Reset()
	}

	return report
}
*/
