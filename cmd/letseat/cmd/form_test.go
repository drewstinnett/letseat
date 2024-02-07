package cmd

import (
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/stretchr/testify/require"
)

func TestNewPlaceOpts(t *testing.T) {
	var target string
	got := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Place").
			Options(newPlaceOpts([]string{"Taco Tuesday"})...).
			Value(&target),
	)).View()

	require.Contains(t, got, "> Someplace New!", "Make sure the default is something new")
	require.Contains(t, got, "Taco Tuesday", "Make sure we still have Taco Tuesday")
}
