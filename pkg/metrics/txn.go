package metrics

import "time"

// Txn is the metrics of one transaction.
type Txn struct {
	Name string // client name

	Err   error
	Size  float64
	Start time.Time
	End   time.Time
}

// TxnSet is collections of Txn, which implements sort.Sort interface.
type TxnSet []Txn

// Append ...
func (s *TxnSet) Append(others []Txn) {
	*s = append(*s, others...)
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
