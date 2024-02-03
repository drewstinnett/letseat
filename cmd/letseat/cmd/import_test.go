package cmd

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
	dbf := path.Join(t.TempDir(), "data.db")
	cmd := newRootCmd()
	cmd.SetArgs([]string{"import", "../testdata/import.yaml", "--data", dbf})
	require.NoError(t, cmd.Execute())
}
