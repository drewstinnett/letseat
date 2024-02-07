package cmd

import (
	"errors"
	"log/slog"

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
	cost     string
	date     string
	takeout  bool
	ratings  map[string]*int
}

func runLog(cmd *cobra.Command, args []string) error {
	diary := letseat.New(
		letseat.WithFilter(*mustNewEntryFilterWithCmd(cmd)),
		letseat.WithDBFilename(mustGetCmd[string](*cmd, "data")),
	)
	e := newEntryForm(nil)

	if err := e.NewForm(diary.Entries()).Run(); err != nil {
		return err
	}

	new := e.Entry()
	gout.MustPrint(new)
	if !doConfirm("Log the entry above?") {
		return errors.New("aborting from confirm, nothing logged")
	}

	if err := diary.Log(new); err != nil {
		return err
	}
	/*
		if err := diary.WriteEntries(); err != nil {
			return err
		}
	*/
	slog.Info("logged!")

	return nil
}

func doConfirm(msg string) bool {
	var confirm bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(msg).
				Value(&confirm),
		),
	).Run(); err != nil {
		return false
	}
	return confirm
}
