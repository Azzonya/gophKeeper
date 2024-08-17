// Package s3 provides methods for interacting with an S3-compatible storage service,
// allowing for file upload, retrieval, and deletion operations.
package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gophKeeper/server/internal/domain/data_items/model"
	"io"
	"log"
	"path/filepath"
	"strconv"
)

// S3Repo manages interactions with the S3 storage, including file operations
// like uploading, retrieving, and deleting objects.
type S3Repo struct {
	client      *minio.Client
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
}

// NewS3Repo initializes a new S3Repo instance, setting up the S3 client and bucket.
// It returns an error if the client creation or bucket setup fails.
func NewS3Repo(ctx context.Context, S3Endpoint, S3AccessKey, S3SecretKey, S3Bucket string) (*S3Repo, error) {
	client, err := minio.New(S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(S3AccessKey, S3SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %v", err)
	}

	err = client.MakeBucket(ctx, S3Bucket, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		exists, errBucketExists := client.BucketExists(ctx, S3Bucket)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", S3Bucket)
		} else {
			log.Fatalf("failed to create bucket: %v", err)
		}
	} else {
		log.Printf("Successfully created bucket %s\n", S3Bucket)
	}

	return &S3Repo{
		client:      client,
		S3Endpoint:  S3Endpoint,
		S3AccessKey: S3AccessKey,
		S3SecretKey: S3SecretKey,
		S3Bucket:    S3Bucket,
	}, nil

}

// GetFile retrieves a file from the S3 bucket based on the provided parameters.
// It returns the file as a byte slice, a boolean indicating if the file exists, and any error encountered.
func (r *S3Repo) GetFile(ctx context.Context, pars *model.GetPars) ([]byte, bool, error) {
	id, _ := strconv.Atoi(pars.ID)
	objectName := filepath.Join("uploads", fmt.Sprintf("%d", id))
	object, err := r.client.GetObject(ctx, r.S3Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, false, fmt.Errorf("failed to get object: %v", err)
	}
	defer object.Close()

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, object)
	if err != nil {
		log.Fatalln(err)
	}

	return buffer.Bytes(), false, nil
}

// UploadFile uploads a file to the S3 bucket, returning the URL of the uploaded file or an error.
func (r *S3Repo) UploadFile(ctx context.Context, id int, data []byte) (string, error) {
	objectName := filepath.Join("uploads", fmt.Sprintf("%d", id))
	_, err := r.client.PutObject(ctx, r.S3Bucket, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %v", err)
	}
	url := fmt.Sprintf("http://%s/%s/%s", r.client.EndpointURL().Host, r.S3Bucket, objectName)

	return url, nil
}

// DeleteFile removes a file from the S3 bucket based on the provided parameters.
// It returns an error if the deletion fails.
func (r *S3Repo) DeleteFile(ctx context.Context, pars *model.GetPars) error {
	objectName := filepath.Join("uploads", pars.ID)
	err := r.client.RemoveObject(ctx, r.S3Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object from MinIO: %v", err)
	}
	return nil
}
