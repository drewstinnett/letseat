package cmd

import (
	"fmt"

	letseat "github.com/drewstinnett/letseat/pkg"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "export entries",
		RunE:  runExport,
	}
	return cmd
}

func runExport(cmd *cobra.Command, args []string) error {
	diary := letseat.New(
		letseat.WithDBFilename(mustGetCmd[string](*cmd, "data")),
	)
	out, err := diary.Export()
	if err != nil {
		return err
	}
	fmt.Fprint(cmd.OutOrStdout(), string(out))
	return diary.Close()
}
