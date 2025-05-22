package minio

import (
	"context"
	"errors"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MinioClient *minio.Client
	endpoint    = "http://127.0.0.1:9000"
	accessKey   = "KiDRzW60OpTZswfxgrvC"
	secretKey   = "hXyf0i2TsgRDR7EDfsrttAqLi4G9Y3ZFST2Y0fFs"
	bucketName  = "skyvault"
	location    = "us-east-1"
	filePath    = "."
)

func InitMinio() {
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
