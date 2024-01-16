package letseat

import (
	"io"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
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

func TestFileWrite(t *testing.T) {
	tmpFile := path.Join(t.TempDir(), "diary.yaml")
	// fh, err := os.Open(tmpFile)
	// require.NoError(t, err)
	err := os.WriteFile(tmpFile, []byte(`- place: Yummy Yums
  date: 2024-01-12`), 0o600)
	require.NoError(t, err)
	d := New(WithEntriesFile(tmpFile))
	d.Log(
		&Entry{
			Place: "mr potatoes bar and grill",
			Date:  toPTR(time.Date(2024, 1, 13, 0, 0, 0, 0, time.UTC)),
		},
	)
	require.NoError(t, d.WriteEntries())

	// Now check the contents
	fh, err := os.Open(tmpFile)
	require.NoError(t, err)
	got, err := io.ReadAll(fh)
	require.NoError(t, err)
	require.Equal(
		t,
		`- place: Yummy Yums
  date: 2024-01-12T00:00:00Z
- place: mr potatoes bar and grill
  date: 2024-01-13T00:00:00Z
`,
		string(got),
	)
}

func TestLog(t *testing.T) {
	d := New()
	d.Log(
		&Entry{
			Place: "heaven",
		},
	)
	require.Equal(
		t,
		&Entries{{Place: "heaven"}},
		d.entries,
	)
	require.Equal(
		t,
		Entries{{Place: "heaven"}},
		d.unfilteredEntries,
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
