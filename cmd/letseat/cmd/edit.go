package cmd

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/huh"
	letseat "github.com/drewstinnett/letseat/pkg"
	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "edit entries",
		RunE:  runEdit,
	}
	return cmd
}

func editEntryOpts(e letseat.Entries) []huh.Option[string] {
	sort.Slice(e, func(i, j int) bool {
		return e[i].Date.After(*e[j].Date)
	})
	ret := make([]huh.Option[string], len(e))
	for idx, entry := range e {
		ret[idx] = huh.NewOption(fmt.Sprintf("%v - %v", entry.Date.Format("2006-01-02"), entry.Place), entry.Key())
	}
	return ret
}

func runEdit(cmd *cobra.Command, args []string) error {
	diary := letseat.New(
		letseat.WithDBFilename(mustGetCmd[string](*cmd, "data")),
	)
	defer dclose(diary)

	var editID string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which entry would you like to edit?").
				Options(
					editEntryOpts(diary.Entries())...,
				).
				Value(&editID),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}

	e, err := diary.Get(editID)
	if err != nil {
		return err
	}

	editForm := newEntryForm(e)
	if err := editForm.NewForm(diary.Entries()).Run(); err != nil {
		return err
	}

	return nil
}
