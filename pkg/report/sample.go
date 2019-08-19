package report

import "time"

// TODO: refactor
// Sample is snapshot of stats
type Sample struct {
	Flag    int64
	Elapsed time.Duration
}

// Samples is collections of Sample, which implments sort.Sort interface.
type Samples []*Sample

// Append ...
func (u *Samples) Append(more []*Sample) {
	*u = append(*u, more...)
}

// TotalFlags  total arritbutes for every transaction.
func (u Samples) TotalFlags() int64 {
	total := int64(0)
	for _, it := range u {
		total += it.Flag
	}
	return total
}

// Total elapsed for all transactions.
func (u Samples) Total() time.Duration {
	total := time.Duration(0)
	for _, it := range u {
		total += it.Elapsed
	}
	return total
}

// Min elapsed
func (u Samples) Min() time.Duration {
	if len(u) == 0 {
		return time.Duration(0)
	}
	return u[0].Elapsed
}

// Max elapsed
func (u Samples) Max() time.Duration {
	if len(u) == 0 {
		return time.Duration(0)
	}
	return u[len(u)-1].Elapsed
}

// At time at percent 90, 99, 99.9 %
func (u Samples) At(percent float32) time.Duration {
	if len(u) == 0 {
		return time.Duration(0)
	}
	id := int(float32(len(u)-1) * percent)
	return u[id].Elapsed
}

// Avg ...
func (u Samples) Avg() time.Duration {
	return u.Total() / time.Duration(len(u))
}

func (u Samples) Len() int {
	return len(u)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (u Samples) Less(i, j int) bool {
	return u[i].Elapsed < u[j].Elapsed
}

// For Sort
// Swap swaps the elements with indexes i and j.
func (u Samples) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
