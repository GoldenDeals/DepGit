package stroage

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorage implements the Storage interface using Minio S3 client
type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

// NewMinioStorage creates a new Minio storage client
func NewMinioStorage(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*MinioStorage, error) {
	// Initialize minio client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Check if bucket exists and create it if it doesn't
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		log.Printf("Created bucket: %s\n", bucketName)
	}

	return &MinioStorage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// Put stores an object in Minio
func (s *MinioStorage) Put(ctx context.Context, namespace string, objname string, obj io.Reader) error {
	// Combine namespace and objname to create the object key
	objectKey := namespace + "/" + objname

	// Check if object already exists
	_, err := s.client.StatObject(ctx, s.bucketName, objectKey, minio.StatObjectOptions{})
	// WARN: Ignoring error
	if err == nil {
		return os.ErrExist
	}

	// Upload the object
	_, err = s.client.PutObject(ctx, s.bucketName, objectKey, obj, -1,
		minio.PutObjectOptions{})
	return err
}

// Get retrieves an object from Minio
func (s *MinioStorage) Get(ctx context.Context, namespace string, objname string) (io.Reader, error) {
	// Combine namespace and objname to create the object key
	objectKey := namespace + "/" + objname

	// Get object
	obj, err := s.client.GetObject(ctx, s.bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	// obj is already an io.Reader, so we can return it directly
	return obj, nil
}

// List returns all objects in a namespace from Minio
func (s *MinioStorage) List(ctx context.Context, namespace string) ([]string, error) {
	var objects []string

	// List all objects with the prefix of the namespace
	objectCh := s.client.ListObjects(ctx, s.bucketName,
		minio.ListObjectsOptions{
			Prefix:    namespace + "/",
			Recursive: true,
		})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		// Remove the namespace/ prefix from the key
		key := object.Key[len(namespace)+1:]
		objects = append(objects, key)
	}

	return objects, nil
}
