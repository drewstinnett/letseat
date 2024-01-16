package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateDate(t *testing.T) {
	require.NoError(t, validateDate("2021-05-30"))
	require.EqualError(t, validateDate("not-a-date"), `parsing time "not-a-date" as "2006-01-02": cannot parse "not-a-date" as "2006"`)
}

func TestValidatePlace(t *testing.T) {
	require.NoError(t, validatePlace("some place"))
	require.EqualError(t, validatePlace(""), "place cannot be empty")
}
