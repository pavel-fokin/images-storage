package imagesstorage

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"

	"image/gif"
	"image/jpeg"
	"image/png"
)

// BBox is a bounding box for image data.
type BBox struct {
	X, Y, W, H int
}

// Valid verifies if a bounding box is valid.
func (bb BBox) Valid() bool {
	return (bb.X >= 0 || bb.Y >= 0) && (bb.W > 0 && bb.H > 0)
}

// CutOut an image by a bounding box.
func (bb BBox) CutOut(imagedata io.Reader) (io.Reader, error) {
	fullimage, formatName, err := image.Decode(imagedata)
	if err != nil {
		return bytes.NewReader([]byte{}), fmt.Errorf("CutOut(): %w", err)
	}

	subimage := fullimage.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(bb.X, bb.Y, bb.W, bb.H))

	buf := new(bytes.Buffer)

	switch formatName {
	case "png":
		err = png.Encode(buf, subimage)
	case "jpeg":
		err = jpeg.Encode(buf, subimage, &jpeg.Options{})
	case "gif":
		err = gif.Encode(buf, subimage, &gif.Options{})
	default:
		err = errors.New("CutOut() unknown format")
	}
	if err != nil {
		return bytes.NewReader([]byte{}), fmt.Errorf("CutOut(): %w", err)
	}
	return buf, nil
}
