package main

import (
	"context"
	"fmt"
	"net"

	appConfig "github.com/adetunjii/streaming-server/config"
	"github.com/adetunjii/streaming-server/internal/service"
	"github.com/adetunjii/streaming-server/ports/awsservice"
	"github.com/adetunjii/streaming-server/ports/grpc"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	// setup application environment
	appConfig, err := appConfig.LoadConfig(".")
	if err != nil {
		logrus.WithContext(ctx).Error("Failed to load application environment variables.")
		panic(err)
	}

	// aws sdk configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		logrus.WithContext(ctx).Error("Failed to load aws config.")
		return
	}

	// setup aws s3 client
	s3Client := awsservice.SetupS3Client(ctx, cfg)

	// setup video service
	vs := service.NewVideoService(ctx, s3Client)

	var port string

	if appConfig.Port == 0 {
		logrus.WithContext(ctx).Warn("GRPC server address not set. Using default value :8080")
		port = ":8080"
	} else {
		port = fmt.Sprintf(":%d", appConfig.Port)
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	logrus.WithContext(ctx).Infof("Listening for connection on %s", port)
	grpc.StartGrpcServer(ctx, listener, vs)
}
