/*
Package letseat is the main thing that decides where to go for dindin
*/
package letseat

import (
	"errors"
	"os"
	"path"
	"slices"
	"sort"
	"time"

	"github.com/montanaflynn/stats"
	"gopkg.in/yaml.v2"
)

// Diary is the thing holding all of your visits and info
type Diary struct {
	unfilteredEntries Entries
	entries           *Entries
	filter            EntryFilter
	fn                string
}

// Entries returns all the entries matching the filter
func (d Diary) Entries() Entries {
	return *d.entries
}

// Log logs a new entry to your diary
func (d *Diary) Log(e *Entry) {
	d.unfilteredEntries = append(d.unfilteredEntries, *e)
	*d.entries = append(*d.entries, *e)
}

// MostPopularPlace just returns the most popular place
func (d Diary) MostPopularPlace() string {
	return mostFrequent(d.entries.placeNames())
}

// PlaceDetails is just some detail summary pieces of the places in your diary
func (d Diary) PlaceDetails() PlaceDetails {
	e := d.Entries()
	places := e.UniquePlaceNames()
	ret := make(PlaceDetails, len(places))
	for idx, place := range places {
		d := e.placeDetails(place)
		ret[idx] = *d
	}
	return ret
}

// WithEntries sets the diary entries on a new Diary object
func WithEntries(e Entries) func(*Diary) {
	return func(d *Diary) {
		d.unfilteredEntries = e
	}
}

// WithFilter sets the entry filter
func WithFilter(f EntryFilter) func(*Diary) {
	return func(d *Diary) {
		d.filter = f
	}
}

// WithEntriesFile adds the entries from a yaml file to a diary
func WithEntriesFile(fn string) func(*Diary) {
	y, err := os.ReadFile(path.Clean(fn))
	if err != nil {
		panic(err)
	}

	var e Entries
	err = yaml.Unmarshal(y, &e)
	if err != nil {
		panic(err)
	}
	return func(d *Diary) {
		d.unfilteredEntries = e
		d.fn = fn
	}
}

// New returns a new Diary object using functional options
func New(opts ...func(*Diary)) *Diary {
	d := &Diary{}
	for _, opt := range opts {
		opt(d)
	}
	d.entries = toPTR(d.unfilteredEntries.filter(&d.filter))
	return d
}

// WriteEntries write the entries back to a yaml file
func (d Diary) WriteEntries() error {
	if d.fn == "" {
		return errors.New("filename cannot be blank")
	}

	entries, err := yaml.Marshal(d.unfilteredEntries)
	if err != nil {
		return err
	}
	sort.SliceStable(d.unfilteredEntries, func(i, j int) bool {
		return !d.unfilteredEntries[i].Date.After(*d.unfilteredEntries[j].Date)
	})

	err = os.WriteFile(d.fn, entries, 0o600)
	if err != nil {
		return err
	}
	return nil
}

// Entries is multiple DiaryEntry objects
type Entries []Entry

// Entry represents a log about your visit to a restaurant
type Entry struct {
	Place     string         `yaml:"place"`
	Cost      int            `yaml:"cost,omitempty"`
	Date      *time.Time     `yaml:"date"`
	IsTakeout bool           `yaml:"takeout,omitempty"`
	Ratings   map[string]int `yaml:"ratings,omitempty"`
}

func (d *Entry) ratingValuesAsFloat64() []float64 {
	ret := make([]float64, len(d.Ratings))
	idx := 0
	for _, v := range d.Ratings {
		ret[idx] = float64(v)
		idx++
	}
	return ret
}

func (d *Entry) averageRating() float64 {
	r := d.ratingValuesAsFloat64()
	if len(r) == 0 {
		return 0
	}
	m, err := stats.Mean(r)
	if err != nil {
		panic(err)
	}
	return m
}

// EntryFilter defiines how to filter a list of entries
type EntryFilter struct {
	Place       string
	OnlyTakeout bool
	OnlyDineIn  bool
	Earliest    *time.Time
	Latest      *time.Time
}

func (e *Entries) people() []string {
	people := []string{}
	for _, entry := range *e {
		for person := range entry.Ratings {
			if !slices.Contains(people, person) {
				people = append(people, person)
			}
		}
	}
	sort.Strings(people)

	return people
}

// PeopleEnhanced returns all the details on people
func (e *Entries) PeopleEnhanced() []Person {
	people := make([]Person, len(e.people()))
	for idx, name := range e.people() {
		p := Person{
			Name:            name,
			PlaceAvgRatings: map[string]float64{},
		}
		ratings := map[string][]int{}
		// Parse through diary ratings
		for _, entry := range *e {
			if entry.Ratings[name] != 0 {
				ratings[entry.Place] = append(ratings[entry.Place], entry.Ratings[name])
			}
		}
		for k, v := range ratings {
			var total float64
			total = 0
			for _, number := range v {
				total += float64(number)
			}
			p.PlaceAvgRatings[k] = total / float64(len(v))
		}
		people[idx] = p
	}
	return people
}

func (e *Entries) placeNames() []string {
	places := make([]string, len(*e))
	for idx, entry := range *e {
		places[idx] = entry.Place
	}

	return places
}

// UniquePlaceNames are just the simple place names as strings
func (e *Entries) UniquePlaceNames() []string {
	places := []string{}
	for _, entry := range *e {
		if !slices.Contains(places, entry.Place) {
			places = append(places, entry.Place)
		}
	}

	sort.Strings(places)
	return places
}

func (e *Entries) placeDetails(place string) *PlaceDetail {
	f := e.filter(&EntryFilter{Place: place})
	dets := &PlaceDetail{
		Name:          place,
		AverageRating: f.averageRating(),
		Visits:        len(f),
	}

	for _, entry := range f {
		if dets.LastVisit == nil || entry.Date.After(*dets.LastVisit) {
			dets.LastVisit = entry.Date
		}
	}
	return dets
}

func (e *Entries) averageRating() float64 {
	var total float64
	for _, entry := range *e {
		total += entry.averageRating()
	}

	return total / float64(len(*e))
}

func (e *Entries) filter(f *EntryFilter) Entries {
	if f == nil {
		return *e
	}

	filtered := Entries{}
	for _, entry := range *e {
		if f.OnlyTakeout && !entry.IsTakeout {
			continue
		}

		if f.OnlyDineIn && entry.IsTakeout {
			continue
		}

		if f.Place != "" && entry.Place != f.Place {
			continue
		}

		if f.Earliest != nil && entry.Date.Before(*f.Earliest) {
			continue
		}

		if f.Latest != nil && entry.Date.After(*f.Latest) {
			continue
		}

		filtered = append(filtered, entry)
	}

	return filtered
}
