package cmd

import (
	"os"

	letseat "github.com/drewstinnett/letseat/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func newImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Args:  cobra.ExactArgs(1),
		Short: "import entries from a flat yaml file",
		RunE:  runImport,
	}
	return cmd
}

func runImport(cmd *cobra.Command, args []string) error {
	diary := letseat.New(
		letseat.WithDBFilename(mustGetCmd[string](*cmd, "data")),
	)
	eb, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	var entries letseat.Entries
	if err := yaml.Unmarshal(eb, &entries); err != nil {
		return err
	}
	for _, entry := range entries {
		entry := entry
		if err := diary.Log(entry); err != nil {
			return err
		}
	}
	return nil
}
