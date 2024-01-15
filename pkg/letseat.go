package letseat

// Foodie is a...a person sort of?
type Foodie struct {
	Places Places  `yaml:"places"`
	Diary  Entries `yaml:"diary"`
}

// AddPlace adds a new place to a foodie
func (f *Foodie) AddPlace(p Place) {
	f.Places = append(f.Places, p)
}

// WithPlaces sets the initial places for a Foodie
func WithPlaces(p Places) func(*Foodie) {
	return func(f *Foodie) {
		f.Places = p
	}
}

// NewFoodie returns a new foodie
func NewFoodie(opts ...func(*Foodie)) (*Foodie, error) {
	f := &Foodie{}
	for _, opt := range opts {
		opt(f)
	}
	return f, nil
}

// MustNewFoodie returns a new foodie or panics on error
func MustNewFoodie(opts ...func(*Foodie)) *Foodie {
	got, err := NewFoodie(opts...)
	if err != nil {
		panic(err)
	}
	return got
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
