package cmd

import (
	"errors"
	"fmt"
	"log/slog"
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
	place    string
	newPlace string
	date     string
	takeout  bool
	ratings  map[string]*int
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

	if e.newPlace != "" {
		ret.Place = e.newPlace
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
	placeOpts := make([]huh.Option[string], len(places)+1)
	placeOpts[0] = huh.Option[string]{
		Key:   "Someplace New!",
		Value: "",
	}
	for idx, item := range places {
		placeOpts[idx+1] = huh.Option[string]{
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

func newRatingInputs(people []letseat.Person, e entryForm) []huh.Field {
	ratingInputs := make([]huh.Field, len(people))
	for idx, item := range people {
		e.ratings[item.Name] = toPTR(0)
		ratingInputs[idx] = huh.NewSelect[int]().
			Title(fmt.Sprintf("%v's Rating", item.Name)).
			Options(ratingOptions...).
			Value(e.ratings[item.Name])
	}
	return ratingInputs
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

	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Date").
				Validate(validateDate).
				Value(&e.date),
			huh.NewSelect[string]().
				Title("Place").
				Options(placeOpts...).
				Value(&e.place),
			huh.NewConfirm().
				Title("Take Out?").
				Value(&e.takeout),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Description("What's this new place called??").
				Validate(validatePlace).
				Value(&e.newPlace),
		).WithHideFunc(func() bool {
			return e.place != ""
		}),
		huh.NewGroup(
			newRatingInputs(entries.PeopleEnhanced(), e)...,
		),
	).Run(); err != nil {
		return err
	}

	new := e.Entry()
	gout.MustPrint(new)
	if !doConfirm() {
		return errors.New("aborting from confirm, nothing logged")
	}

	diary.Log(&new)
	if err := diary.WriteEntries(); err != nil {
		return err
	}
	slog.Info("logged!")

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
