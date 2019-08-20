package stan

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/shohi/nats-bench/pkg/bench"
	"github.com/shohi/nats-bench/pkg/config"
)

// New creates a new `stan` subcommand
func New() *cobra.Command {
	var conf config.StanConfig

	c := &cobra.Command{
		Use:   "stan",
		Short: "benchmark nats-streaming server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bench.RunStanBench(conf)
		},
	}

	setupFlags(c, &conf)

	return c
}

// setupFlags sets flags for comand line
func setupFlags(cmd *cobra.Command, conf *config.StanConfig) {
	fs := cmd.Flags()

	// TODO: add more configurations
	fs.StringVarP(&conf.URL, "server", "s", "nats://localhost:4222", "nats server.")
	fs.DurationVarP(&conf.Duration, "duration", "d", 60*time.Second, "running duration.")
	fs.Float64VarP(&conf.Rate, "rate", "r", 1, "msg publish rate per subject. If set to 0, means no rate limit.")
	fs.StringVarP(&conf.Cluster, "cluster", "c", "nss-cluster", "stan cluster.")

	fs.DurationVar(&conf.PubAckWait, "pat", time.Second*20, "publish ack timeout.")
	fs.IntVar(&conf.PubAckMaxInflight, "mpa", 1000, "max publish ack In flight.")

	fs.IntVar(&conf.PingInterval, "pi", 5, "ping interval.")
	fs.IntVar(&conf.PingMax, "pm", 10, "ping max times without ack.")

	fs.DurationVar(&conf.NatsConnectTimeout, "nct", time.Second*5, "nats connection timeout")
	fs.DurationVar(&conf.StanConnectTimeout, "sct", time.Second*30, "stan connection timeout")

	fs.StringVar(&conf.SubjectPrefix, "sp", "test", "subject prefix.")
	fs.IntVar(&conf.SubjectNum, "sn", 1, "subject number.")
	fs.IntVar(&conf.PubNum, "np", 1, "publisher number per subject.")
	fs.IntVar(&conf.SubNum, "ns", 0, "subscriber number per subject.")
	fs.IntVar(&conf.MsgSize, "ms", 8192, "msg size bytes (8192)")
	fs.IntVarP(&conf.MsgNumber, "count", "n", -1, "msg number per subject. If duration is set, it will be ignored")
	fs.DurationVar(&conf.ReportInterval, "ri", 1*time.Minute, "report interval.")
	fs.BoolVarP(&conf.PubAsync, "async", "a", false, "true for async mode.")

	// cmd.MarkFlagRequired("duration")
	// cmd.MarkFlagRequired("cluster")
}
