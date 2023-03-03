package imagesstorage

import (
	"bytes"
	"fmt"
	"io"

	"github.com/google/uuid"
)

type Image struct {
	UUID        string
	ContentType string
	Width       int
	Height      int
	Size        int
	UploadedAt  string
	Data        []byte
}

func NewImage(uuid uuid.UUID, contenttype string, data io.Reader) (Image, error) {
	buf, err := io.ReadAll(data)
	if err != nil {
		return Image{}, fmt.Errorf("NewImage(): %w", err)
	}

	width, height, err := GetWidthHeight(bytes.NewReader(buf))
	if err != nil {
		return Image{}, fmt.Errorf("NewImage(): %w", err)
	}

	return Image{
		UUID:        uuid.String(),
		ContentType: contenttype,
		Width:       width,
		Height:      height,
		Size:        len(buf),
		Data:        buf,
	}, nil
}
