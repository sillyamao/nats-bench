package report

import "time"

func min(l, r time.Duration) time.Duration {
	if l < r {
		return l
	}
	return r
}

func max(l, r time.Duration) time.Duration {
	if l < r {
		return r
	}

	return l
}

func weightedMedian(t1 time.Duration, c1 int, t2 time.Duration, c2 int) time.Duration {
	return (t1*time.Duration(c1) + t2*time.Duration(c2)) / time.Duration(c1+c2)
}
