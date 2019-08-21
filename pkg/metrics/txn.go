package metrics

import "time"

// Txn is the metrics of one transaction.
type Txn struct {
	Name string // client name

	Err   error
	Size  int64
	Start time.Time
	End   time.Time
}

// TxnSet is collections of Txn, which implements sort.Sort interface.
type TxnSet []*Txn

// Append ...
func (s *TxnSet) Append(others []*Txn) {
	*s = append(*s, others...)
}

// TotalSize returns total transfered size in bytes.
func (s TxnSet) TotalSize() int64 {
	var total int64
	for _, t := range s {
		total += t.Size
	}

	return total
}

// TotalDuration sums up total duration for all transactions.
func (s TxnSet) TotalDuration() time.Duration {
	var total time.Duration
	for _, t := range s {
		total += t.End.Sub(t.Start)
	}

	return total
}

func (s TxnSet) Len() int {
	return len(s)
}

// Less is based on transaction's duration.
func (s TxnSet) Less(i, j int) bool {
	iE := s[i].End.Sub(s[i].Start)
	jE := s[j].End.Sub(s[j].Start)

	return iE < jE
}

// Swap swaps the elements with indexes i and j.
func (s TxnSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

/*
FIXME
// MinDuration
func (s TxnSet) MinDuration() time.Duration {
	if len(s) == 0 {
		return time.Duration(0)
	}

	return u[0].Elapsed
}

// MaxDuration
func (s TxnSet) MaxDuration() time.Duration {
	if len(u) == 0 {
		return time.Duration(0)
	}
	return u[len(u)-1].Elapsed
}

// At time at percent 90, 99, 99.9 %
func (s TxnSet) At(percent float32) time.Duration {
	if len(u) == 0 {
		return time.Duration(0)
	}
	id := int(float32(len(u)-1) * percent)
	return u[id].Elapsed
}

// Avg ...
func (s TxnSet) Avg() time.Duration {
	return u.Total() / time.Duration(len(u))
}
*/
