package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAnalyze(t *testing.T) {
	cmd := &cobra.Command{}
	// bindRootArgs(cmd)
	// bindFilter(cmd)
	runAnalyze(cmd, []string{"--diary ./testdata/bigdiary.yaml"})
}
