package storage

import (
	"context"
	"io"
	"io/ioutil"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"pavel-fokin/images-storage/internal/imagesstorage"
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
	query := &storage.Query{Prefix: ""}
	it := s.bucket.Objects(ctx, query)

	images := []imagesstorage.Image{}

	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("listBucket: unable to list bucket %q: %v", s.config.BucketName, err)
			return []imagesstorage.Image{}, err
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
		log.Println(err)
		return imagesstorage.Image{}, err
	}

	return s.asImage(objAttrs), nil
}

func (s *Storage) Upload(
	ctx context.Context, filename string, contenttype string, data io.Reader,
) (imagesstorage.Image, error) {
	wc := s.bucket.Object(filename).NewWriter(ctx)
	wc.ContentType = contenttype
	defer wc.Close()

	buf, err := ioutil.ReadAll(data)
	if err != nil {
		log.Println(err)
		return imagesstorage.Image{}, err
	}

	wc.Write(buf)

	image := imagesstorage.Image{
		Name:        filename,
		ContentType: contenttype,
		Size:        len(buf),
	}

	return image, nil
}

func (s *Storage) asImage(obj *storage.ObjectAttrs) imagesstorage.Image {
	return imagesstorage.Image{
		Name:        obj.Name,
		ContentType: obj.ContentType,
		Size:        int(obj.Size),
		UploadedAt:  obj.Updated.String(),
	}
}
