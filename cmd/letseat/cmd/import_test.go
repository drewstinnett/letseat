package cmd

import (
	"bytes"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
	dbf := path.Join(t.TempDir(), "data.db")
	cmd := newRootCmd()
	cmd.SetArgs([]string{"import", "../testdata/import.yaml", "--data", dbf})
	require.NoError(t, cmd.Execute())

	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"export", "--data", dbf})
	require.NoError(t, cmd.Execute())
	require.Equal(
		t,
		`- place: Franks Place
  date: 2023-12-14T00:00:00Z
  takeout: true
  ratings:
    andrei: 4
    jeymes: 3
- place: McDonuoughs Pub
  date: 2023-12-16T00:00:00Z
  ratings:
    andrei: 4
- place: Biggy Wings
  date: 2023-12-21T00:00:00Z
  takeout: true
  ratings:
    andrei: 3
`,
		b.String(),
	)
}
