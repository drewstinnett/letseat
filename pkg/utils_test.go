package letseat

import (
	"testing"
	"time"

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

func TestParseDuration(t *testing.T) {
	tests := map[string]struct {
		given     string
		expect    time.Duration
		expectErr string
	}{
		"5 days": {
			given:  "5d",
			expect: time.Hour * 24 * 5,
		},
		"5 days ago": {
			given:  "-5d",
			expect: time.Hour * 24 * 5 * -1,
		},
		"half a day": {
			given:  ".5d",
			expect: time.Hour * 12,
		},
	}
	for desc, tt := range tests {
		if tt.expectErr == "" {
			got, err := ParseDuration(tt.given)
			require.NoError(t, err, desc)
			require.Equal(t, tt.expect, got, desc)
		}
	}
}

func TestStars(t *testing.T) {
	require.Equal(t, "⭐️⭐️⭐️", Stars(1.5, ""))
	require.Equal(t, "🤣🤣🤣", Stars(1.5, "🤣"))
}
