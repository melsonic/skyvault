package minio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
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
		log.Fatal("Please specify a bucket name in .env file")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connected succesfully!")
	MinioClient = client

	err = MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{
		Region: location,
	})

	if err != nil {
		exists, errBucketExists := MinioClient.BucketExists(context.Background(), bucketName)
		if exists && errBucketExists == nil {
			log.Printf("%s already exists\n", bucketName)
		} else {
			log.Fatal(err.Error())
		}
	} else {
		log.Printf("%v created successfully!\n", bucketName)
	}
}

func UploadChunk(hash string, data []byte) error {
	if hash == "" || len(data) == 0 {
		return errors.New("Invalid request data!")
	}
	contentType := "application/octet-stream"
	// saving the chunk with name as 'hash'
	_, err := MinioClient.PutObject(context.Background(), bucketName, hash, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func GetChunk(hash string) ([]byte, error) {
	if hash == "" {
		return nil, errors.New("Invalid request data!")
	}
	object, err := MinioClient.GetObject(context.Background(), bucketName, hash, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	stat, _ := object.Stat()
	result := make([]byte, stat.Size)
	_, err = object.Read(result)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return result, nil
}
