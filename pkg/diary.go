package letseat

import (
	"io/ioutil"
	"sort"
	"time"

	"gopkg.in/yaml.v2"
)

type Diary []DiaryEntry

type DiaryEntry struct {
	Place     string         `yaml:"place"`
	Date      *time.Time     `yaml:"date"`
	IsTakeout bool           `yaml:"takeout"`
	Ratings   map[string]int `yaml:"ratings"`
}

func (d *DiaryEntry) AverageRating() float64 {
	var total float64
	for _, rating := range d.Ratings {
		total += float64(rating)
	}

	return total / float64(len(d.Ratings))
}

type DiaryFilter struct {
	Place       string
	OnlyTakeout bool
	OnlyDineIn  bool
	Earliest    *time.Time
	Latest      *time.Time
}

type PlaceDetail struct {
	Name          string
	AverageRating float64
	LastVisit     *time.Time
	Visits        int
}

func (d *Diary) People() []string {
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

func (d *Diary) PeopleEnhanced() []Person {
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

func (d *Diary) Places() []string {
	places := []string{}
	for _, entry := range *d {
		places = append(places, entry.Place)
	}

	return places
}

func (d *Diary) UniquePlaces() []string {
	places := []string{}
	for _, entry := range *d {
		if !ContainsString(places, entry.Place) {
			places = append(places, entry.Place)
		}
	}

	return places
}

func (d *Diary) PlaceDetails(place string) (*PlaceDetail, error) {
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

func (d *Diary) MostPopularPlace() string {
	places := d.Places()
	return mostFrequent(places)
}

func (d *Diary) AverageRating() float64 {
	var total float64
	for _, entry := range *d {
		total += entry.AverageRating()
	}

	return total / float64(len(*d))
}

func (d *Diary) Filter(f *DiaryFilter) (Diary, error) {
	if f == nil {
		return *d, nil
	}

	filtered := Diary{}
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

func LoadDiaryWithFile(f string, filter *DiaryFilter) (Diary, error) {
	y, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	d := Diary{}
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
