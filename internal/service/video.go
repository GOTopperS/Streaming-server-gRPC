package service

import (
	"bufio"
	"context"
	"fmt"
	"io"

	grpcservice "github.com/adetunjii/streaming-server/pb"
	"github.com/adetunjii/streaming-server/ports/awsservice"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	BufSize = 64 * 1024
)

type VideoService struct {
	s3Client *awsservice.S3Client
}

func NewVideoService(ctx context.Context, s3Client *awsservice.S3Client) *VideoService {
	return &VideoService{
		s3Client: s3Client,
	}
}

/* the goal is to directly read from an s3 bucket and serve it to the client
 * without actually downloading the file from the s3 bucket
 *
 * Basically creating a pipeline to the file in the s3 bucket.
 */
func (v *VideoService) StreamVideo(
	ctx context.Context,
	bucket string,
	videoID string,
	srv grpcservice.VideoService_GetVideoStreamServer,
) error {
	object, err := v.s3Client.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(videoID),
	})

	if err != nil {
		return err
	}

	defer object.Body.Close()

	bufr := bufio.NewReaderSize(object.Body, BufSize)
	part := 1
	for {

		// read data in `Bufsize` byte chunks
		buf := make([]byte, BufSize)
		n, err := io.ReadFull(bufr, buf)
		if err == io.EOF {
			break
		}

		fmt.Printf("Read %d bytes of data \n", n)

		err = srv.Send(&grpcservice.VideoResponse{
			Buffer: buf,
			Part:   int32(part),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
