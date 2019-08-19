package stan

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/shohi/nats-bench/pkg/bench"
	"github.com/shohi/nats-bench/pkg/config"
)

// New creates a new `stan` subcommand
func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "stan",
		Short: "benchmark nats-streaming server",
		RunE:  benchStan,
	}

	return c
}

func benchStan(cmd *cobra.Command, args []string) error {
	var conf config.StanConfig
	setupFlags(cmd, &conf)

	return bench.RunStanBench(conf)
}

// setupFlags sets flags for comand line
func setupFlags(cmd *cobra.Command, conf *config.StanConfig) {
	fs := cmd.Flags()

	// TODO: add more configurations
	fs.Float64VarP(&conf.Rate, "rate", "r", 1, "msg publish rate per topic.")
	fs.IntVar(&conf.SubjectNumber, "sn", 10, "topic number")
	fs.IntVar(&conf.PubNumber, "np", 1, "publisher number per topic")
	fs.IntVar(&conf.SubNumber, "ns", 1, "subscriber number per topic")
	fs.IntVarP(&conf.MsgNumber, "count", "n", -1, "msg number per channel")
	fs.IntVar(&conf.MsgSize, "ms", 8192, "msg size bytes (8192)")
	fs.DurationVarP(&conf.Duration, "duration", "d", 60*time.Second, "running duration")
	fs.DurationVar(&conf.ReportInterval, "ri", 60*time.Second, "report interval")
	fs.BoolVarP(&conf.Async, "async", "a", false, "true for async mode")
	fs.StringVar(&conf.SubjectPrefix, "sp", "test", "subject prefix. `test` is used by default")

	cmd.MarkFlagRequired("duration")
}
