package main

import (
	"time"
)

var (
	layouts = []string{
		"2006",
		"2006-01",
		"2006-01-02",
		"2006-01-02 15",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05 MST",
	}
)

func ParseDateTime(value string) (time.Time, error) {
	var t time.Time
	var e error
	for _, l := range layouts {
		t, e = time.Parse(l, value)
		if e == nil {
			break
		}
	}
	return t, e
}
