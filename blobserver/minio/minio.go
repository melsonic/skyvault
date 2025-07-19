package minio

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MinioClient *minio.Client
	endpoint    string
	accessKey   string
	secretKey   string
	bucketName  string
	location    string
)

func InitMinio() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	endpoint = os.Getenv("MINIO_ENDPOINT")
	accessKey = os.Getenv("ACCESSKEY")
	secretKey = os.Getenv("SECRETKEY")
	bucketName = os.Getenv("BUCKETNAME")
	location = os.Getenv("LOCATION")

	if bucketName == "" {
		slog.Error("bucket name is not specified")
		log.Fatal("Please specify a bucket name in .env file")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}

	slog.Info("Minio connected!")
	MinioClient = client

	exists, errBucketExists := MinioClient.BucketExists(context.Background(), bucketName)

	if exists {
		slog.Info("bucket already exists", "name", bucketName)
		return
	} else if errBucketExists != nil {
		slog.Error("error checking bucket exists", "error", errBucketExists)
	}

	err = MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{
		Region: location,
	})

	if err != nil {
		slog.Error("error creating bucket", "error", err)
		return
	}

	slog.Info("bucket created successfully", "bucket", bucketName)
}

func isObjectExists(hash string) bool {
	_, err := MinioClient.StatObject(context.Background(), bucketName, hash, minio.StatObjectOptions{})
	if err != nil {
		slog.Debug("Object doesn't exist", "ObjectName", hash)
		return false
	}
	slog.Debug("Object already exist", "ObjectName", hash)
	return true
}

func UploadChunk(hash string, data []byte) error {
	if isObjectExists(hash) {
		return nil
	}
	contentType := "application/octet-stream"
	// saving the chunk with name as 'hash'
	_, err := MinioClient.PutObject(context.Background(), bucketName, hash, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		slog.Error(err.Error())
		return errors.New("error uploading object")
	}
	return nil
}

func GetChunk(hash string) ([]byte, error) {
	if isObjectExists(hash) {
		return nil, errors.New("object doesn't exist")
	}
	object, err := MinioClient.GetObject(context.Background(), bucketName, hash, minio.GetObjectOptions{})
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.New("error fetching object")
	}
	stat, _ := object.Stat()
	result := make([]byte, stat.Size)
	_, err = object.Read(result)
	if err != nil && err != io.EOF {
		slog.Error(err.Error())
		return nil, errors.New("error reading object data")
	}
	return result, nil
}

func DeleteChunk(hash string) error {
	err := MinioClient.RemoveObject(context.Background(), bucketName, hash, minio.RemoveObjectOptions{})
	if err != nil {
		slog.Error(err.Error())
		return errors.New("failed to delete object")
	}
	return nil
}
