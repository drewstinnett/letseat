package letseat

import "time"

type Diary []DiaryEntry

type DiaryEntry struct {
	Place     string         `yaml:"place"`
	Date      *time.Time     `yaml:"date"`
	IsTakeout bool           `yaml:"takeout"`
	Ratings   map[string]int `yaml:"ratings"`
}
