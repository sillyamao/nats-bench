package report

import (
	"fmt"
	"time"

	"github.com/shohi/nats-bench/pkg/metrics"
)

// Sample represents statistics of one currency client during one dump interval.
type Sample struct {
	metrics.TxnSet
}

func (s *Sample) Reset() {
	if len(s.TxnSet) > 0 {
		s.TxnSet = s.TxnSet[:0]
	}
}

const (
	reportFormat = `
    Start At:               %v
    Document Size:          %v (KB)
    Concurrency:            %v
    Elapsed:                %v
    Completed:              %v
    Failed:                 %v
    Total Transfered:       %v (MB)
    TPS:                    %v
    Request Time:           %v
    Transfer rate:          %v (MB/s)
`

	formatRecent = `
    Recent:	%v
	Elapsed:				%v
	Completed:				%v
	Failed:					%v
	Total transferred:		%v (MB)
	TPS:					%v
	Request Time:			%v
	Transfer rate:			%v (MB/s)
`
)

// Report content
type Report struct {
	Completed int

	Start         time.Time // report generated time
	Elapsed       time.Duration
	Concurrency   int
	Failed        int
	TotalSize     int64
	TotalDuration time.Duration
	Min, Max, Avg time.Duration
	T90           time.Duration
	T99           time.Duration
	T999          time.Duration
}

func (r *Report) Merge(w Report) {
	// First time, just copy
	if r.Completed == 0 {
		*r = w
		return
	}
	r.T90 = weightedMedian(r.T90, r.Completed, w.T90, w.Completed)
	r.T99 = weightedMedian(r.T99, r.Completed, w.T99, w.Completed)
	r.T999 = weightedMedian(r.T999, r.Completed, w.T999, w.Completed)
	r.TotalDuration += w.TotalDuration
	r.Elapsed += w.Elapsed
	r.Completed += w.Completed
	r.Failed += w.Failed

	r.TotalSize += w.TotalSize
	r.Avg = r.TotalDuration / time.Duration(r.Completed+r.Failed)
	r.Min = min(r.Min, w.Min)
	r.Max = max(r.Max, w.Max)
}

// TODO: double check stats
// String output report's content.
func (r Report) String() string {
	timeStatFmt := fmt.Sprintf("Min:\t%v\tMax:\t%v\tAvg:\t%v\t90:\t%v\t99:\t%v\t99.9:\t%v\t",
		r.Min, r.Max, r.Avg, r.T90, r.T99, r.T999,
	)

	return fmt.Sprintf(reportFormat,
		r.Start,
		float64(r.TotalSize/int64(r.Completed+1))/float64(1<<10),
		r.Concurrency,
		r.Elapsed,
		r.Completed,
		r.Failed,
		float64(r.TotalSize)/float64(1<<20),
		float64(r.Completed)*float64(time.Second)/float64(r.Elapsed),
		timeStatFmt,
		float64(r.TotalSize)*float64(time.Second)/(float64(r.Elapsed)*float64(1<<20)),
	)
}
