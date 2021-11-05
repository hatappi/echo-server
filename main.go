package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

const (
	exitOK = iota
	exitNG
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	ctx := context.Background()

	host := ""
	if h, ok := os.LookupEnv("ECHO_HOST"); ok {
		host = h
	}

	wg, ctx := errgroup.WithContext(ctx)

	httpPort := "3000"
	if p, ok := os.LookupEnv("HTTP_PORT"); ok {
		httpPort = p
	}
	wg.Go(func() error { return runHTTPServer(fmt.Sprintf("%s:%s", host, httpPort)) })

	grpcPort := "5000"
	if p, ok := os.LookupEnv("GRPC_PORT"); ok {
		grpcPort = p
	}
	wg.Go(func() error { return runGRPCServer(fmt.Sprintf("%s:%s", host, grpcPort)) })

	grpcTLSPort := os.Getenv("GRPC_TLS_PORT")
	grpcTLSCrt := os.Getenv("GRPC_TLS_CRT")
	grpcTLSKey := os.Getenv("GRPC_TLS_KEY")
	fmt.Printf("%v\n", grpcTLSPort != "" && grpcTLSCrt != "" && grpcTLSKey != "")
	if grpcTLSPort != "" && grpcTLSCrt != "" && grpcTLSKey != "" {
		wg.Go(func() error { return runTLSGRPCServer(fmt.Sprintf("%s:%s", host, grpcTLSPort), grpcTLSCrt, grpcTLSKey) })
	}

	if err := wg.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] error received: %s\n", err)
		return exitNG
	}

	return exitOK
}
