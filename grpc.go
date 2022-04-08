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

type server struct {
	metadata map[string]string
	logger   logr.Logger
}

func (s *server) SayHello(ctx context.Context, in *pb.SayHelloRequest) (*pb.SayHelloResponse, error) {
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

func runGRPCServer(addr string, metadata map[string]string, logger logr.Logger) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	pb.RegisterEchoServer(s, &server{metadata: metadata, logger: logger})
	reflection.Register(s)

	logger.Info("start gRPC server", "addr", addr)
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
