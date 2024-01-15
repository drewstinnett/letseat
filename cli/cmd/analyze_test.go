package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyze(t *testing.T) {
	b := bytes.NewBufferString("")
	cmd := newRootCmd()
	cmd.SetOut(b)
	cmd.SetArgs([]string{"analyze", "--current-date", "2024-01-12", "--diary", "../testdata/bigdiary.yaml"})
	require.NoError(t, cmd.Execute())
	require.Contains(t, b.String(), "Old Person Wings         62 days ago")
}
