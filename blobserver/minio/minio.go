package minio

import (
	"context"
	"errors"
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
	filePath    string
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
	filePath = os.Getenv("FILEPATH")

	if bucketName == "" {
		log.Fatal("Please specify a bucket name in .env file")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	MinioClient = client

	err = MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{
		Region: location,
	})

	if err != nil {
		exists, errBucketExists := MinioClient.BucketExists(context.Background(), bucketName)
		if exists && errBucketExists != nil {
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
	// saving the chunk with name as 'hash'
	_, err := MinioClient.FPutObject(context.Background(), bucketName, hash, filePath, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
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
	var result []byte
	_, err = object.Read(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
