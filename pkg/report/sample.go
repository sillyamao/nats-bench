package report

import (
	"time"

	"github.com/montanaflynn/stats"
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

type SampleGroup struct {
	Samples []*Sample

	summary Summary

	// durations return duration array which only includes successful Txn.
	// also cast time.Duration to float64 to apply `stats` functions.
	durations []float64

	initialized bool
}

type Summary struct {
	TotalSize    float64
	FailCount    int64
	SuccessCount int64
}

type valueFn func(stats.Float64Data) (float64, error)

// TODO
func (g *SampleGroup) init() {
	if g.initialized {
		return
	}

	ds := make([]float64, 0, len(g.Samples)*1024)
	for _, s := range g.Samples {
		for _, v := range s.TxnSet {
			g.summary.TotalSize += v.Size
			if v.Err != nil {
				g.summary.FailCount++
				continue
			}
			g.summary.SuccessCount++
			ds = append(ds, float64(v.End.Sub(v.Start)))
		}
	}

	g.durations = ds
	g.initialized = true
}

func (g *SampleGroup) Summary() Summary {
	g.init()
	return g.summary
}

func (g *SampleGroup) apply(fn valueFn) time.Duration {
	g.init()

	res, err := fn(g.durations)

	if err != nil {
		return 0
	}

	return time.Duration(res)
}

func (g *SampleGroup) Min() time.Duration {
	return g.apply(stats.Min)
}

func (g *SampleGroup) Max() time.Duration {
	return g.apply(stats.Max)
}

func (g *SampleGroup) Mean() time.Duration {
	return g.apply(stats.Mean)
}

func (g *SampleGroup) Sum() time.Duration {
	return g.apply(stats.Sum)
}

func (g *SampleGroup) Percentile(p float64) time.Duration {
	fn := func(input stats.Float64Data) (float64, error) {
		return stats.Percentile(input, p)
	}

	return g.apply(fn)
}
