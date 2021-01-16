package main

import (
	"testing"
)

// 2021-01-15 2021 "2020-04-05 16:20:02 PDT"
func TestParseDateTime(t *testing.T) {
	values := []string{
		"2021-01-15",
		"2021",
		"2020-04-05 16:20:02 PDT",
		"2021-07-22 03",
	}

	for _, v := range values {
		_, e := ParseDateTime(v)
		if e != nil {
			t.Errorf("%q %v", v, e)
		}
	}
}
