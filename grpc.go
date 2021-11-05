package main

import (
	"context"
	"log"
	"net"

	pb "github.com/hatappi/echo-server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedEchoServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.SayHelloRequest) (*pb.SayHelloReply, error) {
	log.Printf("grpc-server Received: %v", in.GetName())
	return &pb.SayHelloReply{Message: "Hello " + in.GetName()}, nil
}

func runGRPCServer(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	pb.RegisterEchoServer(s, &server{})
	reflection.Register(s)

	log.Printf("start grpc server: addr is %s", addr)
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}

func runTLSGRPCServer(addr, crtPath, keyPath string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	creds, err := credentials.NewServerTLSFromFile(crtPath, keyPath)
	if err != nil {
		return err
	}

	// Create the gRPC server with the credentials
	s := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterEchoServer(s, &server{})
	reflection.Register(s)

	log.Printf("start TLS grpc server: addr is %s", addr)
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
