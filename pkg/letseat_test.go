package letseat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFoodie(t *testing.T) {
	got, err := New()
	require.NoError(t, err)
	require.NotNil(t, got)
}

func TestAddPlaces(t *testing.T) {
	f, err := New(WithPlaces(Places{
		{Name: "foo"},
	}))
	require.NoError(t, err)
	require.NotNil(t, f)

	f.AddPlace(Place{
		Name: "bar",
	})
	require.Equal(
		t,
		Foodie{
			Places: Places{
				{
					Name: "Foo",
					Slug: "foo",
				},
				{
					Name: "bar",
					Slug: "bar",
				},
			},
		},
		*f,
	)
}
