package imagesstorage

import (
	"bytes"
	"image"
	"io"
	"log"

	_ "image/gif"
	_ "image/jpeg"
	"image/png"
)

type BBox struct {
	X, Y, W, H int
}

func (bb BBox) Valid() bool {
	return (bb.X > 0 || bb.Y > 0) || (bb.W > 0 && bb.H > 0)
}

func GetWidthHeight(imagedata io.Reader) (width int, height int) {
	img, _, err := image.Decode(imagedata)
	if err != nil {
		log.Println(err)
		return 0, 0
	}
	return img.Bounds().Max.X, img.Bounds().Max.Y
}

func CutOut(imagedata io.Reader, bbox BBox) io.Reader {
	fullimage, _, _ := image.Decode(imagedata)

	subimage := fullimage.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(bbox.X, bbox.Y, bbox.W, bbox.H))

	buf := new(bytes.Buffer)
	// TODO only png now
	png.Encode(buf, subimage)
	return buf
}
