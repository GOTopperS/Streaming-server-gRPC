package grpc

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/adetunjii/streaming-server/internal/service"
	grpcservice "github.com/adetunjii/streaming-server/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	grpcservice.UnimplementedVideoServiceServer
	grpcServer   *grpc.Server
	listener     net.Listener
	videoService *service.VideoService
}

func (s GrpcServer) Close() error {
	s.grpcServer.GracefulStop()
	return nil
}

func StartGrpcServer(ctx context.Context, listener net.Listener, vs *service.VideoService) {

	srv := GrpcServer{
		listener:     listener,
		videoService: vs,
	}

	srv.grpcServer = grpc.NewServer()
	grpcservice.RegisterVideoServiceServer(srv.grpcServer, &srv)

	logrus.WithContext(ctx).Info("Grpc Server started successfully...")
	if err := srv.grpcServer.Serve(srv.listener); err != nil {
		logrus.WithContext(ctx).Fatal(fmt.Sprintf("Cannot start grpc server: %s", err))
	}

}

func (g *GrpcServer) GetVideoStream(req *grpcservice.VideoRequest, srv grpcservice.VideoService_GetVideoStreamServer) error {
	ctx := context.TODO()

	logrus.WithContext(ctx).Infof("Request received with %s", req.Id)

	err := g.videoService.StreamVideo(ctx, "gotsstreamingserver", req.Id, srv)
	if err != nil {
		return err
	}

	return nil
}

var (
	_ io.Closer                      = (*GrpcServer)(nil)
	_ grpcservice.VideoServiceServer = (*GrpcServer)(nil)
)
