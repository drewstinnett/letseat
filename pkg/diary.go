/*
Package letseat is the main thing that decides where to go for dindin
*/
package letseat

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"sort"
	"time"

	"github.com/montanaflynn/stats"
	bolt "go.etcd.io/bbolt"
)

// Diary is the thing holding all of your visits and info
type Diary struct {
	// Deprecated: Use db instead of unfiltered entries
	unfilteredEntries Entries
	entries           *Entries
	filter            EntryFilter
	db                *bolt.DB
}

// Entries returns all the entries matching the filter
func (d Diary) Entries() Entries {
	return *d.entries
}

// Log logs a new entry to your diary
func (d *Diary) Log(es ...Entry) error {
	if d.entries == nil {
		d.entries = &Entries{}
	}
	if err := d.db.Update(func(tx *bolt.Tx) error {
		for _, e := range es {
			e := e
			if err := tx.Bucket([]byte(EntriesBucket)).Put([]byte(e.Key()), e.mustMarshal()); err != nil {
				slog.Warn("error logging entry", "error", err)
			}
			*d.entries = append(*d.entries, e)
		}
		// Now save the person info
		for _, person := range d.entries.PeopleEnhanced() {
			if err := tx.Bucket([]byte(PeopleBucket)).Put([]byte(person.Name), []byte("true")); err != nil {
				slog.Warn("error logging person", "error", err)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
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

// WithDB sets the bbolt database for a letseat client
func WithDB(db *bolt.DB) func(*Diary) {
	if err := initDB(db); err != nil {
		panic(err)
	}
	var entries Entries
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(EntriesBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var entry Entry
			if err := json.Unmarshal(v, &entry); err != nil {
				panic(err)
			}
			entries = append(entries, entry)
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return func(d *Diary) {
		d.db = db
		d.unfilteredEntries = entries
	}
}

// WithDBFilename uses a given file for the db
func WithDBFilename(fn string) func(*Diary) {
	db, err := bolt.Open(fn, 0o600, nil)
	if err != nil {
		panic(err)
	}
	return WithDB(db)
}

// WithEntries sets the diary entries on a new Diary object
func WithEntries(e Entries) func(*Diary) {
	return func(d *Diary) {
		if d.db == nil {
			panic("must set the db before loading any entries")
		}
		ents := make([]Entry, len(e))
		for idx, ent := range e {
			ent := ent
			ents[idx] = ent
		}
		if err := d.Log(ents...); err != nil {
			panic(err)
		}
		d.unfilteredEntries = e
	}
}

// WithFilter sets the entry filter
func WithFilter(f EntryFilter) func(*Diary) {
	return func(d *Diary) {
		d.filter = f
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

const (
	// PeopleBucket is the name of the bucket that contains all the People entries
	PeopleBucket = "people"
	// EntriesBucket is the name of the bucket that contains all the Entries entries
	EntriesBucket = "entries"
	// PlacesBucket is the name of the bucket that contains all the Places entries
	PlacesBucket = "places"
)

func initDB(db *bolt.DB) error {
	buckets := []string{PeopleBucket, EntriesBucket, PlacesBucket}
	for _, bucket := range buckets {
		bucket := bucket
		if err := db.Update(func(tx *bolt.Tx) error {
			if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

/*
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
*/

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

// Key is the key path for the database for a given entry
func (d Entry) Key() string {
	if d.Date == nil {
		d.Date = &time.Time{}
	}
	return fmt.Sprintf("/%v/%v", d.Date.Format(time.RFC3339), d.Place)
}

func (d Entry) mustMarshal() []byte {
	got, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return got
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

	return people
}

// PeopleEnhanced returns all the details on people
func (e *Entries) PeopleEnhanced() []Person {
	if len(e.people()) == 0 {
		return []Person{}
	}
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
