package letseat

import "sort"

type Person struct {
	Name            string
	PlaceAvgRatings map[string]float64
}

func (p *Person) FavoriteN(n int) []string {
	places := []string{}
	type kv struct {
		Key   string
		Value float64
	}
	var ss []kv
	for k, v := range p.PlaceAvgRatings {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	for i, p := range ss {
		if i >= n {
			break
		}
		places = append(places, p.Key)
	}
	return places
}
