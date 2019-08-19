package bench

// Runner represents a benchmark routine
type Runner interface {
	Run()
	Stop()
}
