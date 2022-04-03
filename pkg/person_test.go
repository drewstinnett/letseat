package letseat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPeople(t *testing.T) {
	ts := []struct {
		diary Diary
		want  []string
	}{
		{
			diary: Diary{
				DiaryEntry{
					Place: "A", Ratings: map[string]int{
						"a": 1,
						"b": 2,
					},
				},
				DiaryEntry{
					Place: "A", Ratings: map[string]int{
						"c": 1,
						"d": 2,
					},
				},
			},
			want: []string{"a", "b", "c", "d"},
		},
	}

	for _, tt := range ts {
		got := tt.diary.People()
		require.Equal(t, tt.want, got)
	}
}

func TestFavoriteN(t *testing.T) {
	ts := []struct {
		diary Diary
		n     int
		want  []string
	}{
		{
			diary: Diary{
				DiaryEntry{Place: "yum", Ratings: map[string]int{"a": 5, "b": 4}},
				DiaryEntry{Place: "yuck", Ratings: map[string]int{"a": 1, "b": 2}},
			},
			n:    2,
			want: []string{"yum", "yuck"},
		},
		{
			diary: Diary{
				DiaryEntry{Place: "yuck", Ratings: map[string]int{"a": 1, "b": 2}},
				DiaryEntry{Place: "yum", Ratings: map[string]int{"a": 5, "b": 4}},
			},
			n:    2,
			want: []string{"yum", "yuck"},
		},
		{
			diary: Diary{
				DiaryEntry{Place: "yuck", Ratings: map[string]int{"a": 1, "b": 2}},
				DiaryEntry{Place: "yum", Ratings: map[string]int{"a": 5, "b": 4}},
			},
			n:    1,
			want: []string{"yum"},
		},
	}

	for _, tt := range ts {
		people := tt.diary.PeopleEnhanced()
		for _, person := range people {
			if person.Name == "a" {
				got := person.FavoriteN(tt.n)
				require.Equal(t, tt.want, got)
			}
		}
	}
}
