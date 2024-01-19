package cmd

import (
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "configuration",
		Aliases: []string{"config", "c"},
		Short:   "Print configuration paths",
		RunE:    runConfig,
	}
	return cmd
}

func runConfig(cmd *cobra.Command, args []string) error {
	gout.MustPrint(config)

	return nil
}
