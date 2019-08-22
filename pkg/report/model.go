package report

import (
	"fmt"
	"time"
)

const reportFormat = `
    StartTime:         %v
    Elapsed:           %v
    Concurrency:       %v
    Success/Failure:   %v
    TPS:               %v
    AvgRespSize:       %v (KB)
    RespTimeDist:      %v
    TotalTransfered:   %v (MB)
    TransferRate:      %v (MB/s)
`

// Report content
type Report struct {
	Start       time.Time // report generated time
	Elapsed     time.Duration
	Concurrency int

	Success   int64
	Failed    int64
	TotalSize float64

	Total time.Duration
	Min   time.Duration
	Max   time.Duration
	Mean  time.Duration
	T90   time.Duration
	T99   time.Duration
	T999  time.Duration
}

// Populate uses sample group to fill report.
func (r *Report) Populate(sg *SampleGroup) {
	sa := sg.Summary()

	r.Success = sa.SuccessCount
	r.Failed = sa.FailCount
	r.TotalSize = sa.TotalSize

	r.Total = sg.Sum()
	r.Min = sg.Min()
	r.Max = sg.Max()
	r.Mean = sg.Mean()
	r.T90 = sg.Percentile(90)
	r.T99 = sg.Percentile(99)
	r.T999 = sg.Percentile(99.9)
}

func (r *Report) Merge(t Report) {
	if (r.Success + r.Failed) == 0 {
		*r = t
		return
	}

	r.Success += t.Success
	r.Failed += t.Failed
	r.TotalSize += t.TotalSize

	r.T90 = weightedMedian(r.T90, r.Success, t.T90, t.Success)
	r.T99 = weightedMedian(r.T99, r.Success, t.T99, t.Success)
	r.T999 = weightedMedian(r.T999, r.Success, t.T999, t.Success)
	r.Total += t.Total
	r.Elapsed += t.Elapsed

	r.Mean = r.Total / time.Duration(r.Success+r.Failed)
	r.Min = min(r.Min, t.Min)
	r.Max = max(r.Max, t.Max)
}

// TODO: double check stats
// String output report's content.
func (r Report) String() string {
	dDist := fmt.Sprintf("[Min:%v Max:%v Mean:%v 90:%v 99:%v 99.9:%v]",
		r.Min, r.Max, r.Mean, r.T90, r.T99, r.T999,
	)

	succFail := fmt.Sprintf("%v/%v", r.Success, r.Failed)

	// in KB
	var avgRespSize float64
	if r.Success > 0 {
		avgRespSize = (r.TotalSize / float64(r.Success)) / float64(1<<10)
	}

	return fmt.Sprintf(reportFormat,
		r.Start,
		r.Elapsed,
		r.Concurrency,
		succFail,
		float64(r.Success)/r.Elapsed.Seconds(),
		avgRespSize,
		dDist,
		float64(r.TotalSize)/float64(1<<20),
		float64(r.TotalSize)/(r.Elapsed.Seconds()*float64(1<<20)),
	)
}
