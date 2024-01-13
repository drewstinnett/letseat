package letseat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEntryAverage(t *testing.T) {
	ts := []struct {
		entry Entry
		want  float64
	}{
		{
			Entry{
				Place: "a",
				Ratings: map[string]int{
					"a": 1,
					"b": 2,
					"c": 3,
				},
			},
			2,
		},
	}

	for _, tt := range ts {
		got := tt.entry.averageRating()
		require.Equal(t, tt.want, got)
	}
}

func TestFilterDiary(t *testing.T) {
	ts := []struct {
		diary  Entries
		filter *DiaryFilter
		want   int
	}{
		{
			diary: Entries{
				Entry{Place: "A", IsTakeout: true},
				Entry{Place: "B", IsTakeout: false},
			},
			filter: &DiaryFilter{
				OnlyTakeout: true,
			},
			want: 1,
		},
		{
			diary: Entries{
				Entry{Place: "A", IsTakeout: true},
				Entry{Place: "B", IsTakeout: false},
			},
			filter: &DiaryFilter{
				OnlyDineIn: true,
			},
			want: 1,
		},
		{
			diary: Entries{
				Entry{Place: "A", IsTakeout: true},
				Entry{Place: "B", IsTakeout: false},
			},
			filter: &DiaryFilter{
				Place: "A",
			},
			want: 1,
		},
		{
			diary: Entries{
				Entry{Place: "A", IsTakeout: true},
				Entry{Place: "B", IsTakeout: false},
			},
			want: 2,
		},
	}
	for _, tt := range ts {
		got, err := tt.diary.Filter(tt.filter)
		require.NoError(t, err)
		require.Equal(t, tt.want, len(got))
	}
}

func TestPopularPlace(t *testing.T) {
	ts := []struct {
		diary Entries
		want  string
		msg   string
	}{
		{
			diary: Entries{
				Entry{Place: "B"},
				Entry{Place: "A"},
				Entry{Place: "C"},
				Entry{Place: "A"},
			},
			want: "A",
			msg:  "Short list of entries",
		},
		{
			diary: Entries{
				Entry{Place: "A"},
				Entry{Place: "B"},
				Entry{Place: "C"},
				Entry{Place: "A"},
				Entry{Place: "A"},
				Entry{Place: "B"},
				Entry{Place: "C"},
			},
			want: "A",
			msg:  "Longer list of entries",
		},
	}
	for _, tt := range ts {
		got := tt.diary.MostPopularPlace()
		require.Equal(t, tt.want, got)
	}
}
