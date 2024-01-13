/*
Package letseat is the main thing that decides where to go for dindin
*/
package letseat

import (
	"os"
	"sort"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Diary is the thing holding all of your visits and info
type Diary struct {
	Entries Entries
}

// Entries is multiple DiaryEntry objects
type Entries []Entry

// Entry represents a log about your visit to a restaurant
type Entry struct {
	Place     string         `yaml:"place"`
	Cost      int            `yaml:"cost"`
	Date      *time.Time     `yaml:"date"`
	IsTakeout bool           `yaml:"takeout"`
	Ratings   map[string]int `yaml:"ratings"`
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

type DiaryFilter struct {
	Place       string
	OnlyTakeout bool
	OnlyDineIn  bool
	Earliest    *time.Time
	Latest      *time.Time
}

func NewDiaryFilterWithCmd(cmd *cobra.Command) (*DiaryFilter, error) {
	onlyTakeout, err := cmd.Flags().GetBool("only-takeout")
	if err != nil {
		return nil, err
	}
	onlyDinein, err := cmd.Flags().GetBool("only-dinein")
	if err != nil {
		return nil, err
	}

	// Earliest
	earliestA, err := cmd.Flags().GetString("earliest")
	if err != nil {
		return nil, err
	}
	earliestD, err := ParseDuration(earliestA)
	if err != nil {
		return nil, err
	}
	earliest := time.Now().Add(-earliestD)
	return &DiaryFilter{
		OnlyTakeout: onlyTakeout,
		OnlyDineIn:  onlyDinein,
		Earliest:    &earliest,
	}, nil
}

func (d *Entries) People() []string {
	people := []string{}
	for _, entry := range *d {
		for person := range entry.Ratings {
			if !ContainsString(people, person) {
				people = append(people, person)
			}
		}
	}
	sort.Strings(people)

	return people
}

func (d *Entries) PeopleEnhanced() []Person {
	people := []Person{}
	for _, name := range d.People() {
		ratings := map[string][]int{}
		p := Person{
			Name:            name,
			PlaceAvgRatings: map[string]float64{},
		}
		// Parse through diary ratings
		for _, entry := range *d {
			if entry.Ratings[name] != 0 {
				ratings[entry.Place] = append(ratings[entry.Place], entry.Ratings[name])
			}
		}
		for k, v := range ratings {
			var total float64
			total = 0
			for _, number := range v {
				total = total + float64(number)
			}
			avg := total / float64(len(v))
			p.PlaceAvgRatings[k] = avg
		}
		people = append(people, p)
	}
	return people
}

func (d *Entries) Places() []string {
	places := []string{}
	for _, entry := range *d {
		places = append(places, entry.Place)
	}

	return places
}

func (d *Entries) UniquePlaces() []string {
	places := []string{}
	for _, entry := range *d {
		if !ContainsString(places, entry.Place) {
			places = append(places, entry.Place)
		}
	}

	return places
}

func (d *Entries) PlaceDetails(place string) (*PlaceDetail, error) {
	e, err := d.Filter(&DiaryFilter{Place: place})
	if err != nil {
		return nil, err
	}
	dets := &PlaceDetail{
		Name:          place,
		AverageRating: e.AverageRating(),
		Visits:        len(e),
	}

	for _, entry := range e {
		if dets.LastVisit == nil || entry.Date.After(*dets.LastVisit) {
			dets.LastVisit = entry.Date
		}
	}
	return dets, nil
}

func (d *Entries) MostPopularPlace() string {
	places := d.Places()
	return mostFrequent(places)
}

func (d *Entries) AverageRating() float64 {
	var total float64
	for _, entry := range *d {
		total += entry.averageRating()
	}

	return total / float64(len(*d))
}

func (d *Entries) Filter(f *DiaryFilter) (Entries, error) {
	if f == nil {
		return *d, nil
	}

	filtered := Entries{}
	for _, entry := range *d {
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

	return filtered, nil
}

func LoadDiaryWithFile(f string, filter *DiaryFilter) (Entries, error) {
	y, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	d := Entries{}
	err = yaml.Unmarshal(y, &d)
	if err != nil {
		return nil, err
	}
	d, err = d.Filter(filter)
	if err != nil {
		return nil, err
	}

	return d, nil
}
