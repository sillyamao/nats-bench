package report

import (
	"context"
	"fmt"
	"sort"
	"time"
)

// TODO: refactor
type Reports struct {
	Times       int
	TotalReport Report
	Start       time.Time
	Lables      map[string]*History
}

func (r *Reports) Get(label string) *History {
	if r.Lables == nil {
		r.Start = time.Now()
		r.Lables = make(map[string]*History, 4096)
	}
	his := r.Lables[label]
	if his == nil {
		his = &History{
			Errors:  0,
			Samples: make([]*Sample, 1<<16)[:0],
		}
		r.Lables[label] = his
	}
	return his
}

// Report for period
type Report struct {
	At            time.Time
	Elapsed       time.Duration
	Concurrency   int
	Failed        int
	Completed     int
	TotalFlags    int64
	Total         time.Duration
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
	r.Total += w.Total
	r.Elapsed += w.Elapsed
	r.Completed += w.Completed
	r.Failed += w.Failed
	r.TotalFlags += w.TotalFlags
	r.Avg = r.Total / time.Duration(r.Completed)
	r.Min = min(r.Min, w.Min)
	r.Max = max(r.Max, w.Max)
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

// Print the report
func (r Report) Print() string {
	timestat := fmt.Sprintf("Min:\t%v\tMax:\t%v\tAvg:\t%v\t90:\t%v\t99:\t%v\t99.9:\t%v\t",
		r.Min, r.Max, r.Avg, r.T90, r.T99, r.T999,
	)

	return fmt.Sprintf(reportFormat,
		r.At,
		float64(r.TotalFlags/int64(r.Completed+1))/1000.0,
		r.Concurrency,
		r.Elapsed,
		r.Completed,
		r.Failed,
		float64(r.TotalFlags)/float64(1<<20),
		float64(r.Completed)*float64(time.Second)/float64(r.Elapsed),
		timestat,
		// r.ElapsedTotal/time.Duration(r.Completed+1),
		float64(r.TotalFlags)*float64(time.Second)/(float64(r.Elapsed)*float64(1<<20)),
	)
}

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

// Build ... build report data and return print string.
func (r *Reports) Build() string {
	r.Times++
	report := Report{}
	report.At = r.Start
	report.Elapsed = time.Since(report.At)
	report.Concurrency = len(r.Lables)
	labels := make([]string, len(r.Lables))[:0]
	var us Samples
	for key, value := range r.Lables {
		labels = append(labels, key)
		report.Failed += value.Errors
		report.Completed += len(value.Samples)
		us.Append(value.Samples)
	}
	report.Report(us)
	r.TotalReport.Merge(report)
	return fmt.Sprintf(`
Total:    %v times, at:%v
%v
Recent:%v
	`,
		r.Times,
		time.Now(),
		r.TotalReport.Print(), report.Print())
}

// BuildTotal build total report
func (r *Reports) BuildTotal() string {
	report := Report{}
	report.At = r.Start
	report.Elapsed = time.Since(report.At)
	report.Concurrency = len(r.Lables)
	labels := make([]string, len(r.Lables))[:0]
	var us Samples
	for key, value := range r.Lables {
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
	return fmt.Sprintf(`
Reports:%v
	`, r.TotalReport.Print())
}

// BuildAndClear and clear history
func (r *Reports) BuildAndClear() string {
	report := r.Build()
	r.Start = time.Now()
	for _, value := range r.Lables {
		value.Reset()
	}
	return report
}

// Reporter ...
type Reporter struct {
	PrintInterval time.Duration
	ReportTo      string
	Reports       Reports
	Operation     chan func()
	Context       context.Context
}

// NewReporter ...
func NewReporter(ctx context.Context, reportInterval time.Duration) *Reporter {
	reporter := &Reporter{}
	reporter.PrintInterval = reportInterval
	reporter.Context = ctx
	reporter.Operation = make(chan func(), 128)
	return reporter
}

// Start ...
func (r *Reporter) Start() {
	r.Reports.Start = time.Now()

	go func() {
		ticker := time.NewTicker(r.PrintInterval)
		defer func() {
			fmt.Println("reporter finished")
		}()
		defer ticker.Stop()
		for {
			select {
			case <-r.Context.Done():
				return
			case op := <-r.Operation:
				op()
			case <-ticker.C:
				fmt.Println(r.Reports.BuildAndClear())
				// fmt.Printf("Report at:\t%v\n     From Last Report:\n%v\n",
				// 	time.Now(),
				// 	// r.GReports.Build(),
				// 	r.Reports.BuildAndClear(),
				// )
			}
		}
	}()
}

func (r *Reporter) error(label string) {
	// r.GReports.Get(label).Error()
	r.Reports.Get(label).Error()
}

func (r *Reporter) add(label string, size int, elapsed time.Duration) {
	unit := &Sample{int64(size), elapsed}
	// r.GReports.Get(label).Add(unit)
	r.Reports.Get(label).Add(unit)
}

// func (r *Reporter) GetGReport() string {
// 	ch := make(chan string)
// 	r.Operation <- func() {
// 		ch <- r.GReports.Build()
// 	}
// 	return <-ch
// }

// GetReport from reporter.
func (r *Reporter) GetReport() string {
	ch := make(chan string)
	r.Operation <- func() {
		ch <- r.Reports.BuildTotal()
	}
	return <-ch
}

// Report a sample
func (r *Reporter) Report(err error, label string, size int, elapsed time.Duration) {
	r.Operation <- func() {
		if err != nil {
			r.error(label)
		} else {
			r.add(label, size, elapsed)
		}
	}
}
