package letseat

import (
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

func newTestDB(t *testing.T) *bolt.DB {
	dbf := path.Join(t.TempDir(), "test.db")
	db, err := bolt.Open(dbf, 0o600, nil)
	require.NoError(t, err)
	return db
}

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
		d := New(
			WithDB(newTestDB(t)),
			WithEntries(tt.entries),
		)
		got := d.MostPopularPlace()
		require.Equal(t, tt.want, got)
	}
}

func TestNew(t *testing.T) {
	require.NotNil(t, New(
		WithDB(newTestDB(t)),
	))
	d := New(
		WithDB(newTestDB(t)),
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

func TestLog(t *testing.T) {
	d := New(
		WithDB(newTestDB(t)),
	)
	d.Log(
		Entry{
			Place: "heaven",
		},
	)
	require.Equal(
		t,
		&Entries{{Place: "heaven"}},
		d.entries,
	)
}

func TestEntryUnmarshal(t *testing.T) {
	y := `place: Mamacitas
date: 2024-01-15T00:00:00Z
takeout: true
ratings:
  drew: 5
  james: 3`

	var got Entry
	require.NoError(t, yaml.Unmarshal([]byte(y), &got))
	require.Equal(
		t,
		Entry{
			Place:     "Mamacitas",
			Date:      toPTR(time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)),
			IsTakeout: true, Ratings: map[string]int{
				"drew":  5,
				"james": 3,
			},
		},
		got,
	)
}

func TestWithDB(t *testing.T) {
	db := newTestDB(t)
	got := New(WithDB(db))
	require.NotNil(t, got)
	require.NotNil(t, got.db)
	require.NotNil(t, New(WithDBFilename(path.Join(t.TempDir(), "test-fn.db"))))
}
