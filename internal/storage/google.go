package storage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type Config struct {
	BucketName string `env:"IMAGES_STORAGE_GOOGLE_BUCKET_NAME" envDefault:""`
}

type Storage struct {
	config Config
	client *storage.Client
	bucket *storage.BucketHandle
}

func New(ctx context.Context, config Config) *Storage {
	client, err := storage.NewClient(ctx)
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

func (s *Storage) List(ctx context.Context) error {

	query := &storage.Query{Prefix: ""}
	it := s.bucket.Objects(ctx, query)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("listBucket: unable to list bucket %q: %v", s.config.BucketName, err)
			return err
		}

		fmt.Printf("%s %s\n", obj.Name, obj.ContentType)
	}

	return nil
}

func (s *Storage) Upload(
	ctx context.Context, filename string, contenttype string, data io.Reader,
) error {
	wc := s.bucket.Object(filename).NewWriter(ctx)
	wc.ContentType = contenttype
	defer wc.Close()

	buf, err := ioutil.ReadAll(data)
	if err != nil {
		log.Println(err)
		return err
	}

	wc.Write(buf)

	return nil
}
