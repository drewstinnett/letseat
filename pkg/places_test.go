package letseat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPlace(t *testing.T) {
	got, err := NewPlace(WithName("Foo Bar"), WithFormat(Format{
		DineIn: true,
	}))
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(
		t,
		&Place{
			Name: "Foo Bar",
			Slug: "foo-bar",
			Format: Format{
				DineIn: true,
			},
		},
		got,
	)
}

func TestNewPlaceEmptyName(t *testing.T) {
	got, err := NewPlace()
	require.Nil(t, got)
	require.EqualError(t, err, "name cannot be empty")
}
