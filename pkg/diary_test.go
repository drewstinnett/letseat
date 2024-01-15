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
			entry: Entry{
				Place: "a",
				Ratings: map[string]int{
					"a": 1,
					"b": 2,
					"c": 3,
				},
			},
			want: 2,
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
		filter *EntryFilter
		want   int
	}{
		{
			diary: Entries{
				Entry{Place: "A", IsTakeout: true},
				Entry{Place: "B", IsTakeout: false},
			},
			filter: &EntryFilter{
				OnlyTakeout: true,
			},
			want: 1,
		},
		{
			diary: Entries{
				Entry{Place: "A", IsTakeout: true},
				Entry{Place: "B", IsTakeout: false},
			},
			filter: &EntryFilter{
				OnlyDineIn: true,
			},
			want: 1,
		},
		{
			diary: Entries{
				Entry{Place: "A", IsTakeout: true},
				Entry{Place: "B", IsTakeout: false},
			},
			filter: &EntryFilter{
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
		got := tt.diary.filter(tt.filter)
		require.Equal(t, tt.want, len(got))
	}
}

func TestPopularPlace(t *testing.T) {
	ts := []struct {
		entries Entries
		want    string
		msg     string
	}{
		{
			entries: Entries{
				Entry{Place: "B"},
				Entry{Place: "A"},
				Entry{Place: "C"},
				Entry{Place: "A"},
			},
			want: "A",
			msg:  "Short list of entries",
		},
		{
			entries: Entries{
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
		d := New(WithEntries(tt.entries))
		got := d.MostPopularPlace()
		require.Equal(t, tt.want, got)
	}
}

func TestNew(t *testing.T) {
	require.NotNil(t, New())
	d := New(
		WithEntries(
			Entries{
				Entry{Place: "Some Dine-In Place"},
				Entry{Place: "Some Takeout Place", IsTakeout: true},
			},
		),
		WithFilter(
			EntryFilter{
				OnlyTakeout: true,
			},
		),
	)
	require.NotNil(t, d)
	require.Equal(
		t,
		Entries{
			Entry{Place: "Some Takeout Place", IsTakeout: true},
		},
		d.Entries(),
	)
}
