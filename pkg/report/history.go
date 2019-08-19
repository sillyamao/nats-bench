package report

// TODO: refactor
// History ...
type History struct {
	Errors  int
	Samples []*Sample
}

// Error ...
func (h *History) Error() {
	h.Errors++
}

// Reset ...
func (h *History) Reset() {
	h.Errors = 0
	h.Samples = h.Samples[:0]
}

// Add ...
func (h *History) Add(u *Sample) {
	h.Samples = append(h.Samples, u)
}
