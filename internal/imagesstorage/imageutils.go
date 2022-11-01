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

type BBox struct {
	X, Y, W, H int
}

func (bb BBox) Valid() bool {
	return (bb.X > 0 || bb.Y > 0) || (bb.W > 0 && bb.H > 0)
}

func GetWidthHeight(imagedata io.Reader) (width int, height int, err error) {
	img, _, err := image.Decode(imagedata)
	if err != nil {
		return 0, 0, fmt.Errorf("GetWidthHeight(): %w", err)
	}
	return img.Bounds().Max.X, img.Bounds().Max.Y, nil
}

func CutOut(imagedata io.Reader, bbox BBox) (io.Reader, error) {
	fullimage, formatName, err := image.Decode(imagedata)
	if err != nil {
		return bytes.NewReader([]byte{}), fmt.Errorf("CutOut(): %w", err)
	}

	subimage := fullimage.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(bbox.X, bbox.Y, bbox.W, bbox.H))

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
