package imagesstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BBoxValid(t *testing.T) {
	tests := []struct {
		bbox     BBox
		expected bool
	}{
		{BBox{0, 0, 0, 0}, false},
		{BBox{0, 0, 10, 0}, false},
		{BBox{0, 0, 0, 10}, false},
		{BBox{10, 10, 0, 0}, false},
		{BBox{0, 0, 10, 10}, true},
		{BBox{10, 0, 10, 10}, true},
		{BBox{10, 10, 10, 10}, true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.bbox.Valid())
	}
}
