package letseat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFoodie(t *testing.T) {
	require.NotNil(t, newFoodie())
}

func TestAddPlaces(t *testing.T) {
	f := newFoodie(withPlaces(Places{*MustNewPlace(WithName("Foo"))}))
	require.NotNil(t, f)
	require.EqualValues(t,
		foodie{
			Places: Places{
				{Name: "Foo", Slug: "foo"},
			},
		},
		*f,
	)

	p := MustNewPlace(WithName("bar"))
	f.addPlace(*p)

	require.Equal(
		t,
		foodie{
			Places: Places{
				{Name: "Foo", Slug: "foo"},
				{Name: "bar", Slug: "bar"},
			},
		},
		*f,
	)
}
