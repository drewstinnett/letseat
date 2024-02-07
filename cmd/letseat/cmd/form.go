package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	letseat "github.com/drewstinnett/letseat/pkg"
)

func (e entryForm) Entry() letseat.Entry {
	d, err := time.Parse("2006-01-02", e.date)
	panicIfErr(err)

	cost, err := strconv.Atoi(e.cost)
	panicIfErr(err)

	ret := letseat.Entry{
		Place:     e.place,
		Date:      &d,
		IsTakeout: e.takeout,
		Ratings:   make(map[string]int, len(e.ratings)),
		Cost:      cost,
	}

	if e.newPlace != "" {
		ret.Place = e.newPlace
	}
	for person, rating := range e.ratings {
		ret.Ratings[person] = *rating
	}
	return ret
}

func (e *entryForm) NewForm(entries letseat.Entries) *huh.Form {
	placeOpts := newPlaceOpts(entries.UniquePlaceNames())

	groups := []*huh.Group{
		huh.NewGroup(
			huh.NewInput().
				Title("Date").
				Description("When did you go?").
				Validate(validateDate).
				Value(&e.date),
			huh.NewSelect[string]().
				Title("Place").
				Description("What's this place called?").
				Options(placeOpts...).
				Value(&e.place),
			huh.NewInput().
				Title("Cost").
				Description("Use 0 for unknown cost").
				Placeholder("0").
				Validate(validateNumber).
				Prompt("$ ").
				Value(&e.cost),
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
	}
	ri := e.newRatingInputs(entries.PeopleEnhanced())
	if len(ri) > 0 {
		groups = append(groups, huh.NewGroup(ri...))
	}
	return huh.NewForm(groups...)
}

var ratingOptions []*huh.Option[int] = []*huh.Option[int]{
	{Key: "🚫 No Rating", Value: 0},
	{Key: "⭐️⭐️⭐️⭐️⭐️", Value: 5},
	{Key: "⭐️⭐️⭐️⭐️", Value: 4},
	{Key: "⭐️⭐️⭐️", Value: 3},
	{Key: "⭐️⭐️", Value: 2},
	{Key: "⭐️", Value: 1},
}

func ratingOptionsWithSelected(s int) []huh.Option[int] {
	var ret []huh.Option[int]
	for _, item := range ratingOptions {
		if item.Value == s {
			ret = append(ret, item.Selected(true))
		} else {
			ret = append(ret, *item)
		}
	}
	return ret
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

// return a new EntryForm using a given entry as a template
func newEntryForm(t *letseat.Entry) entryForm {
	if t == nil {
		return entryForm{
			date:    time.Now().Format("2006-01-02"),
			cost:    "0",
			ratings: map[string]*int{},
		}
	}
	ratings := map[string]*int{}
	for k, v := range t.Ratings {
		v := v
		ratings[k] = &v
	}
	return entryForm{
		date:    t.Date.Format("2006-01-02"),
		cost:    fmt.Sprint(t.Cost),
		ratings: ratings,
		place:   t.Place,
		takeout: t.IsTakeout,
	}
}

func (e *entryForm) newRatingInputs(people []letseat.Person) []huh.Field {
	ratingInputs := make([]huh.Field, len(people))
	for idx, item := range people {
		e.ratings[item.Name] = toPTR(0)
		ro := ratingOptionsWithSelected(*e.ratings[item.Name])
		ratingInputs[idx] = huh.NewSelect[int]().
			Title(fmt.Sprintf("%v's Rating (%v)", item.Name, *e.ratings[item.Name])).
			Options(ro...).
			Value(e.ratings[item.Name])
	}
	return ratingInputs
}
