package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/shohi/nats-bench/cmd/nats"
	"github.com/shohi/nats-bench/cmd/stan"
)

// NewNATSBenchCommand creates the `nats-bench` command and its nested children.
func NewNATSBenchCommand() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "nats-bench",
		Short: "nats benchmark tools",
		Run:   runHelp,
	}

	cmds.AddCommand(nats.New())
	cmds.AddCommand(stan.New())

	return cmds
}

// Execute executes commands, which should be called by main.main() once.
func Execute() {
	command := NewNATSBenchCommand()

	if err := command.Execute(); err != nil {
		log.Printf("run nats-bench error, err: %v\n", err)
		os.Exit(1)
	}
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
