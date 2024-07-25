package grpc

import (
	"context"
	"fmt"
	"io"
	"net"

	grpcservice "github.com/adetunjii/streaming-server/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	grpcservice.UnimplementedHelloServiceServer
	grpcServer *grpc.Server
	listener net.Listener
}

func (s *GrpcServer) Close() error {
	s.grpcServer.GracefulStop()
	return nil
}

func StartGrpcServer(ctx context.Context, addr string) (io.Closer, error) {
	
	if addr == "" {
		logrus.WithContext(ctx).Warn("GRPC server address not set. Using default value :8000")
		addr = ":8000"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	srv := GrpcServer {
		listener: listener,
	}

	srv.grpcServer = grpc.NewServer()
	grpcservice.RegisterHelloServiceServer(srv.grpcServer, &srv)
	
	go func() {
		logrus.WithContext(ctx).Info(fmt.Sprintf("GRPC server is listening on: %s", addr))

		if err = srv.grpcServer.Serve(srv.listener); err != nil {
			logrus.WithContext(ctx).Fatal(fmt.Sprintf("Cannot start grpc server: %s", err))
		}
	}()

	return &srv, nil
}

var (
	_ io.Closer = (*GrpcServer)(nil)
)