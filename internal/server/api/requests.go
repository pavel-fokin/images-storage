package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/pavel-fokin/images-storage/internal/imagesstorage"
)

func ParseBBoxParam(r *http.Request) (imagesstorage.BBox, error) {
	bboxParam := r.URL.Query().Get("bbox")

	if bboxParam == "" {
		return imagesstorage.BBox{}, nil
	}

	params := strings.Split(bboxParam, ",")

	if len(params) != 4 {
		return imagesstorage.BBox{}, ErrBBoxValidate
	}

	bbox := imagesstorage.BBox{}
	bbox.X, _ = strconv.Atoi(params[0])
	bbox.Y, _ = strconv.Atoi(params[1])
	bbox.W, _ = strconv.Atoi(params[2])
	bbox.H, _ = strconv.Atoi(params[3])

	if !bbox.Valid() {
		return imagesstorage.BBox{}, ErrBBoxValidate
	}

	return bbox, nil
}
