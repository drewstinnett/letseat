package cmd

/*
func TestAnalyze(t *testing.T) {
	b := bytes.NewBufferString("")
	cmd := newRootCmd()
	cmd.SetOut(b)
	dbf := path.Join(t.TempDir(), "data.db")

	cmd.SetArgs([]string{"import", "../testdata/bigdiary.yaml", "--data", dbf})
	require.NoError(t, cmd.Execute())

	cmd.SetArgs([]string{"analyze", "--data", dbf})
	require.NoError(t, cmd.Execute())
	assert.Contains(t, b.String(), "Old Person Wings         62 days ago")
}
*/
