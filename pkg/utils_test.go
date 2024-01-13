package letseat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMostFrequent(t *testing.T) {
	tt := []struct {
		values []string
		want   string
	}{
		{[]string{"a", "b", "a", "a", "b", "c", "a", "b", "c"}, "a"},
	}

	for _, ts := range tt {
		require.Equal(t, ts.want, mostFrequent(ts.values))
	}
}