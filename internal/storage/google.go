package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/pavel-fokin/images-storage/internal/imagesstorage"
)

type Config struct {
	BucketName string `env:"IMAGES_STORAGE_GOOGLE_BUCKET_NAME,notEmpty" envDefault:""`
}

type Storage struct {
	config Config
	client *storage.Client
	bucket *storage.BucketHandle
}

func New(config Config) *Storage {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	bucket := client.Bucket(config.BucketName)

	return &Storage{
		config: config,
		client: client,
		bucket: bucket,
	}
}

func (s *Storage) List(ctx context.Context) ([]imagesstorage.Image, error) {
	it := s.bucket.Objects(ctx, &storage.Query{})

	images := []imagesstorage.Image{}

	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []imagesstorage.Image{}, fmt.Errorf("List(): %w", err)
		}

		images = append(images, s.asImage(obj))
	}

	return images, nil
}

func (s *Storage) Metadata(
	ctx context.Context, filename string,
) (imagesstorage.Image, error) {
	objAttrs, err := s.bucket.Object(filename).Attrs(ctx)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrObjectNotExist):
			return imagesstorage.Image{}, imagesstorage.ErrImageNotExist
		default:
			return imagesstorage.Image{}, fmt.Errorf("Metadata(): %w", err)
		}
	}

	return s.asImage(objAttrs), nil
}

func (s *Storage) DoesExist(
	ctx context.Context, filename string,
) bool {
	_, err := s.bucket.Object(filename).Attrs(ctx)
	switch {
	case errors.Is(err, storage.ErrObjectNotExist):
		return true
	default:
		return false
	}
}

func (s *Storage) Download(
	ctx context.Context, uuid string,
) (io.Reader, string, error) {
	r, err := s.bucket.Object(uuid).NewReader(ctx)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrObjectNotExist):
			return bytes.NewReader([]byte{}), "", imagesstorage.ErrImageNotExist
		default:
			return bytes.NewReader([]byte{}), "", fmt.Errorf("Download(): %w", err)
		}
	}
	defer r.Close()

	return r, r.Attrs.ContentType, nil
}

func (s *Storage) Upload(
	ctx context.Context, filename string, contenttype string, data io.Reader, metadata map[string]string,
) (imagesstorage.Image, error) {
	wc := s.bucket.Object(filename).NewWriter(ctx)
	wc.ContentType = contenttype
	wc.Metadata = metadata
	defer wc.Close()

	buf, err := io.ReadAll(data)
	if err != nil {
		return imagesstorage.Image{}, fmt.Errorf("Upload(): %w", err)
	}
	_, err = wc.Write(buf)
	if err != nil {
		return imagesstorage.Image{}, fmt.Errorf("Upload(): %w", err)
	}

	width, _ := strconv.Atoi(metadata["ImageWidth"])
	height, _ := strconv.Atoi(metadata["ImageHeight"])
	image := imagesstorage.Image{
		UUID:        filename,
		ContentType: contenttype,
		Width:       width,
		Height:      height,
		Size:        len(buf),
	}

	return image, nil
}

func (s *Storage) asImage(obj *storage.ObjectAttrs) imagesstorage.Image {
	width, _ := strconv.Atoi(obj.Metadata["ImageWidth"])
	height, _ := strconv.Atoi(obj.Metadata["ImageHeight"])
	return imagesstorage.Image{
		UUID:        obj.Name,
		ContentType: obj.ContentType,
		Size:        int(obj.Size),
		UploadedAt:  obj.Updated.String(),
		Width:       width,
		Height:      height,
	}
}
