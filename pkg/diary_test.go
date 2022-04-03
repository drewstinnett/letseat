package letseat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEntryAverage(t *testing.T) {
	ts := []struct {
		entry DiaryEntry
		want  float64
	}{
		{
			DiaryEntry{
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
		got := tt.entry.AverageRating()
		require.Equal(t, tt.want, got)
	}
}

func TestFilterDiary(t *testing.T) {
	ts := []struct {
		diary  Diary
		filter *DiaryFilter
		want   int
	}{
		{
			diary: Diary{
				DiaryEntry{Place: "A", IsTakeout: true},
				DiaryEntry{Place: "B", IsTakeout: false},
			},
			filter: &DiaryFilter{
				OnlyTakeout: true,
			},
			want: 1,
		},
		{
			diary: Diary{
				DiaryEntry{Place: "A", IsTakeout: true},
				DiaryEntry{Place: "B", IsTakeout: false},
			},
			filter: &DiaryFilter{
				OnlyDineIn: true,
			},
			want: 1,
		},
		{
			diary: Diary{
				DiaryEntry{Place: "A", IsTakeout: true},
				DiaryEntry{Place: "B", IsTakeout: false},
			},
			filter: &DiaryFilter{
				Place: "A",
			},
			want: 1,
		},
		{
			diary: Diary{
				DiaryEntry{Place: "A", IsTakeout: true},
				DiaryEntry{Place: "B", IsTakeout: false},
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
		diary Diary
		want  string
		msg   string
	}{
		{
			diary: Diary{
				DiaryEntry{Place: "B"},
				DiaryEntry{Place: "A"},
				DiaryEntry{Place: "C"},
				DiaryEntry{Place: "A"},
			},
			want: "A",
			msg:  "Short list of entries",
		},
		{
			diary: Diary{
				DiaryEntry{Place: "A"},
				DiaryEntry{Place: "B"},
				DiaryEntry{Place: "C"},
				DiaryEntry{Place: "A"},
				DiaryEntry{Place: "A"},
				DiaryEntry{Place: "B"},
				DiaryEntry{Place: "C"},
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
