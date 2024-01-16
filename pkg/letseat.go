package letseat

// foodie is a...a person sort of?
type foodie struct {
	Places Places  `yaml:"places"`
	Diary  Entries `yaml:"diary"`
}

// addPlace adds a new place to a foodie
func (f *foodie) addPlace(p Place) {
	f.Places = append(f.Places, p)
}

// withPlaces sets the initial places for a Foodie
func withPlaces(p Places) func(*foodie) {
	return func(f *foodie) {
		f.Places = p
	}
}

// newFoodie returns a new foodie
func newFoodie(opts ...func(*foodie)) *foodie {
	f := &foodie{}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

/*
// NewPlace returns a new place object
func NewPlace(s string) (*Place, error) {
	if s == "" {
		return errors.New("name must not be empty")
	}
	p := Place{
		Name: s,
	}
}
*/
