package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
	"io"
	"log"
	"time"
)

type BucketBasics struct {
	S3Client *s3.Client
	Bucket   string
	Key      string
}

func NewS3(conf *viper.Viper) *BucketBasics {

	client := s3.New(s3.Options{
		Region: conf.GetString("aws.region"),
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			conf.GetString("aws.access_key"),
			conf.GetString("aws.secret_key"),
			""))})

	return &BucketBasics{
		S3Client: client,
		Bucket:   conf.GetString("aws.bucket_name"),
	}
}

func (basics BucketBasics) UploadFile(fileName string, body io.Reader) (url string, err error) {

	datePath := time.Now().Format("2006/01/02")
	fileKey := fmt.Sprintf("%v/%v", datePath, fileName)

	fileUrl := fmt.Sprintf("https://track5music.s3.ap-east-1.amazonaws.com/%v", fileKey)

	_, err = basics.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(basics.Bucket),
		Key:    &fileKey,
		Body:   body,
	})

	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			fileName, basics.Bucket, fileKey, err)
		return "", err
	}

	return fileUrl, nil
}
