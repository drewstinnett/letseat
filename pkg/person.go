package letseat

import "sort"

// Person represents a person who ate and rated something at a restaurant
type Person struct {
	Name            string
	PlaceAvgRatings map[string]float64
}

// FavoriteN returns the persons N favorite restaurants
func (p *Person) FavoriteN(n int) []string {
	type kv struct {
		Key   string
		Value float64
	}
	var ss []kv
	for k, v := range p.PlaceAvgRatings {
		ss = append(ss, kv{Key: k, Value: v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	places := []string{}
	for i, p := range ss {
		if i >= n {
			break
		}
		places = append(places, p.Key)
	}
	return places
}
