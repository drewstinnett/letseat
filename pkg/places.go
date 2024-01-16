package letseat

import (
	"errors"
	"time"

	"github.com/gosimple/slug"
)

// Place is a restaurant, or place you can eat
type Place struct {
	Name   string `yaml:"name"`
	Slug   string `yaml:"slug"`
	Tier   int    `yaml:"tier"`
	Format Format `yaml:"format"`
}

// PlaceDetail is the overview detail thing of a place
type PlaceDetail struct {
	Name          string
	AverageRating float64
	LastVisit     *time.Time
	Visits        int
}

// PlaceDetails represents multiple PlaceDetail items. Satisfies the Sortable interface
type PlaceDetails []PlaceDetail

// Len shows the size of the PlaceDetails
func (p PlaceDetails) Len() int {
	return len(p)
}

// Less returns if items are less than other items. Satisfies the Sortable interface
func (p PlaceDetails) Less(i, j int) bool {
	return p[i].AverageRating > p[j].AverageRating
}

// Swap swaps 2 items in PlaceDetails. Satisfies the Sortable interface
func (p PlaceDetails) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// WithName sets the name of a place using functional options
func WithName(n string) func(*Place) {
	return func(p *Place) {
		p.Name = n
		p.Slug = slug.Make(p.Name)
	}
}

// WithFormat sets the format of a place using functional options
func WithFormat(f Format) func(*Place) {
	return func(p *Place) {
		p.Format = f
	}
}

// NewPlace uses functional options to return a new *Place and an optional error
func NewPlace(options ...func(*Place)) (*Place, error) {
	p := &Place{}
	for _, opt := range options {
		opt(p)
	}
	if p.Name == "" {
		return nil, errors.New("name cannot be empty")
	}
	return p, nil
}

// MustNewPlace returns a new place or panics on error
func MustNewPlace(options ...func(*Place)) *Place {
	got, err := NewPlace(options...)
	if err != nil {
		panic(err)
	}
	return got
}

// Places represents multiple Place objects
type Places []Place

// Format represents the type of ways you can dine at at given Place
type Format struct {
	DineIn    bool
	TakeOut   bool
	FoodTruck bool
	Counter   bool
}
