package nats

import (
	"github.com/spf13/cobra"
)

// New creates a new `stan` subcommand
func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "nats",
		Short: "benchmark nats server",
		RunE:  benchNats,
	}

	return c
}

// TODO
func benchNats(cmd *cobra.Command, args []string) error {
	return nil
}
