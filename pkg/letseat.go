package letseat

type Foodie struct {
	Places Places  `yaml:"places"`
	Diary  Entries `yaml:"diary"`
}

func (f *Foodie) AddPlace(p Place) {
	// np := make([]Place, len(f.Places)+1)
	// np[len(f.Places)] = p
	f.Places = append(f.Places, p)
}

func WithPlaces(p Places) func(*Foodie) {
	return func(f *Foodie) {
		f.Places = p
	}
}

func New(options ...func(*Foodie)) (*Foodie, error) {
	f := &Foodie{}
	return f, nil
}
