package imagesstorage

import (
	"image"
	"io"
	"log"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func GetWidthHeight(imagedata io.Reader) (width int, height int) {
	img, _, err := image.Decode(imagedata)
	if err != nil {
		log.Println(err)
		return 0, 0
	}
	return img.Bounds().Max.X, img.Bounds().Max.Y
}

func CutOut(imagedata io.Reader) {

}
