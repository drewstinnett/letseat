package letseat

import (
	"errors"
	"time"

	"github.com/gosimple/slug"
)

type Place struct {
	Name   string `yaml:"name"`
	Slug   string `yaml:"slug"`
	Tier   int    `yaml:"tier"`
	Format Format `yaml:"format"`
}

type PlaceDetail struct {
	Name          string
	AverageRating float64
	LastVisit     *time.Time
	Visits        int
}

func WithName(n string) func(*Place) {
	return func(p *Place) {
		p.Name = n
		p.Slug = slug.Make(p.Name)
	}
}

func WithFormat(f Format) func(*Place) {
	return func(p *Place) {
		p.Format = f
	}
}

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

type Places []Place

type Format struct {
	DineIn    bool
	TakeOut   bool
	FoodTruck bool
	Counter   bool
}
