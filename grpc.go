package main

import (
	"context"
	"net"
	"strings"

	"github.com/go-logr/logr"
	pb "github.com/hatappi/echo-server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type echoServer struct {
	metadata map[string]string
	logger   logr.Logger
}

func (s *echoServer) SayHello(ctx context.Context, in *pb.SayHelloRequest) (*pb.SayHelloResponse, error) {
	reqMeta := make(map[string]string)

	meta, _ := metadata.FromIncomingContext(ctx)
	for k, vals := range meta {
		reqMeta[k] = strings.Join(vals, ",")
	}

	s.logger.Info("sayHello request received", "name", in.GetName(), "protocol", "gRPC")

	return &pb.SayHelloResponse{
		Body:             "Hello " + in.GetName(),
		RequestMetadata:  reqMeta,
		ResponseMetadata: s.metadata,
	}, nil
}

type gRPCServer struct {
	server *grpc.Server
	logger logr.Logger
}

func NewGRPCServer(metadata map[string]string, logger logr.Logger) *gRPCServer {
	s := grpc.NewServer()

	pb.RegisterEchoServer(s, &echoServer{metadata: metadata, logger: logger})
	reflection.Register(s)

	return &gRPCServer{
		server: s,
		logger: logger,
	}
}

func (gs *gRPCServer) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	gs.logger.Info("start gRPC server", "addr", addr)
	if err := gs.server.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (gs *gRPCServer) Shutdown() {
	gs.server.GracefulStop()
}
