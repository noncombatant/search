package main

import (
	"testing"
)

func TestParseSize(t *testing.T) {
	v, e := ParseSize("42Gi")
	if e != nil {
		t.Error(e)
	}
	if v != 42*1024*1024*1024 {
		t.Error(v)
	}
}
