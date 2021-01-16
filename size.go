package main

import (
	"strconv"
	"strings"
)

func ParseSize(value string) (int64, error) {
	scale := int64(1)
	if strings.HasSuffix(value, "T") {
		value = value[:len(value)-1]
		scale = 1000 * 1000 * 1000 * 1000
	} else if strings.HasSuffix(value, "Ti") {
		value = value[:len(value)-2]
		scale = 1024 * 1024 * 1024 * 1024
	} else if strings.HasSuffix(value, "G") {
		value = value[:len(value)-1]
		scale = 1000 * 1000 * 1000
	} else if strings.HasSuffix(value, "Gi") {
		value = value[:len(value)-2]
		scale = 1024 * 1024 * 1024
	} else if strings.HasSuffix(value, "M") {
		value = value[:len(value)-1]
		scale = 1000 * 1000
	} else if strings.HasSuffix(value, "Mi") {
		value = value[:len(value)-2]
		scale = 1024 * 1024
	} else if strings.HasSuffix(value, "K") {
		value = value[:len(value)-1]
		scale = 1000
	} else if strings.HasSuffix(value, "Ki") {
		value = value[:len(value)-2]
		scale = 1024
	}

	n, e := strconv.ParseInt(strings.TrimSpace(value), 0, 64)
	if e != nil {
		return 0, e
	}
	return n * scale, nil
}
