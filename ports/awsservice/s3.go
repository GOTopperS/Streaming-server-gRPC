package awsservice

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

type S3Client struct {
	Client *s3.Client
}

func SetupS3Client(ctx context.Context, config aws.Config) *S3Client {

	return &S3Client{
		Client: s3.NewFromConfig(config),
	}
}

func (s *S3Client) DownloadObject(ctx context.Context, bucket string, key string, filename string) error {
	object, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to get object %s:: %v", key, err)
		return err
	}

	defer object.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to create file %s:: %v", filename, err)
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(object.Body)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to read object content:: %v", err)
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to write content to file %v", filename)
		return err
	}

	return nil
}
