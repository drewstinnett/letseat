package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/drewstinnett/gout/v2"
	letseat "github.com/drewstinnett/letseat/pkg"
	"github.com/spf13/cobra"
)

func newLogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "log a visit to a restaurant or takeout experience",
		RunE:  runLog,
	}
	bindFilter(cmd)
	return cmd
}

type entryForm struct {
	place   string
	date    string
	takeout bool
	ratings map[string]*int
}

func (e entryForm) Entry() letseat.Entry {
	d, err := time.Parse("2006-01-02", e.date)
	panicIfErr(err)
	ret := letseat.Entry{
		Place:     e.place,
		Date:      &d,
		IsTakeout: e.takeout,
		Ratings:   make(map[string]int, len(e.ratings)),
	}
	for person, rating := range e.ratings {
		ret.Ratings[person] = *rating
	}
	return ret
}

var ratingOptions []huh.Option[int] = []huh.Option[int]{
	{Key: "⭐️⭐️⭐️⭐️⭐️", Value: 5},
	{Key: "⭐️⭐️⭐️⭐️", Value: 4},
	{Key: "⭐️⭐️⭐️", Value: 3},
	{Key: "⭐️⭐️", Value: 2},
	{Key: "⭐️", Value: 1},
}

func newPlaceOpts(places []string) []huh.Option[string] {
	placeOpts := make([]huh.Option[string], len(places))
	for idx, item := range places {
		placeOpts[idx] = huh.Option[string]{
			Key:   item,
			Value: item,
		}
	}
	return placeOpts
}

func newEntryForm() entryForm {
	e := entryForm{
		date:    time.Now().Format("2006-01-02"),
		ratings: map[string]*int{},
	}
	return e
}

func runLog(cmd *cobra.Command, args []string) error {
	diary := letseat.New(
		letseat.WithFilter(*mustNewEntryFilterWithCmd(cmd)),
		letseat.WithEntriesFile(mustGetCmd[string](*cmd, "diary")),
	)
	entries := diary.Entries()
	// diary.Entries()
	placeOpts := newPlaceOpts(entries.UniquePlaceNames())

	// Set up the form
	e := newEntryForm()

	// Get the people here
	people := entries.PeopleEnhanced()
	ratingInputs := make([]huh.Field, len(people))
	for idx, item := range people {
		e.ratings[item.Name] = toPTR(0)
		ratingInputs[idx] = huh.NewSelect[int]().
			Description(fmt.Sprintf("%v's Rating", item.Name)).
			Options(ratingOptions...).
			Value(e.ratings[item.Name])
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Description("Date").
				Validate(validateDate).
				Value(&e.date),
			huh.NewSelect[string]().
				Description("Place").
				Options(placeOpts...).
				Validate(validatePlace).
				Value(&e.place),
			huh.NewConfirm().
				Description("Take Out?").
				Value(&e.takeout),
		),
		huh.NewGroup(
			ratingInputs...,
		),
	)
	if err := form.Run(); err != nil {
		return err
	}

	new := e.Entry()
	gout.MustPrint(new)
	if !doConfirm() {
		return errors.New("aborting from confirm, nothing logged")
	}

	diary.Log(&new)

	gout.MustPrint(diary.Entries())

	return nil
}

func doConfirm() bool {
	var confirm bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Description("Create the new entry above?").
				Value(&confirm),
		),
	).Run(); err != nil {
		return false
	}
	return confirm
}
