package api

import (
	"strconv"
	"strings"
)

func ParseBBoxParam(bbox string) (x, y, w, h int) {
	params := strings.Split(bbox, ",")
	x, _ = strconv.Atoi(params[0])
	y, _ = strconv.Atoi(params[1])
	w, _ = strconv.Atoi(params[2])
	h, _ = strconv.Atoi(params[3])
	return x, y, w, h
}
