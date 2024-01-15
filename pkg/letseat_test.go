package letseat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFoodie(t *testing.T) {
	got, err := NewFoodie()
	require.NoError(t, err)
	require.NotNil(t, got)
}

func TestAddPlaces(t *testing.T) {
	f := MustNewFoodie(WithPlaces(Places{*MustNewPlace(WithName("Foo"))}))
	require.NotNil(t, f)
	require.EqualValues(t,
		Foodie{
			Places: Places{
				{Name: "Foo", Slug: "foo"},
			},
		},
		*f,
	)

	p := MustNewPlace(WithName("bar"))
	f.AddPlace(*p)

	require.Equal(
		t,
		Foodie{
			Places: Places{
				{Name: "Foo", Slug: "foo"},
				{Name: "bar", Slug: "bar"},
			},
		},
		*f,
	)
}
