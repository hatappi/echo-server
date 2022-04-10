package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/hatappi/go-kit/log/zap"
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
	logger, err := zap.NewLogger("echo-server")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger. %s\n", err)
		return exitNG
	}

	meta := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], "ECHO_SERVER_META_") {
			meta[strings.TrimPrefix(pair[0], "ECHO_SERVER_META_")] = pair[1]
		}
	}

	host := ""
	if h, ok := os.LookupEnv("ECHO_SERVER_HOST"); ok {
		host = h
	}

	ctx := context.Background()
	wg, ctx := errgroup.WithContext(ctx)

	httpServer := NewHTTPServer(meta, logger)

	httpPort := "3000"
	if p, ok := os.LookupEnv("ECHO_SERVER_HTTP_PORT"); ok {
		httpPort = p
	}
	wg.Go(func() error { return httpServer.Run(fmt.Sprintf("%s:%s", host, httpPort)) })

	grpcServer := NewGRPCServer(meta, logger)

	grpcPort := "5000"
	if p, ok := os.LookupEnv("ECHO_SERVER_GRPC_PORT"); ok {
		grpcPort = p
	}
	wg.Go(func() error { return grpcServer.Run(fmt.Sprintf("%s:%s", host, grpcPort)) })

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	<-ctx.Done()

	var shutdownWG sync.WaitGroup

	shutdownWG.Add(2)

	go func() {
		logger.Info("start shutdown HTTP server")
		// need to pass new context due to existing context is already done
		if err := httpServer.Shutdown(context.Background()); err != nil {
			logger.Error(err, "failed to shutdown HTTP server")
		}
		logger.Info("end shutdown HTTP server")

		shutdownWG.Done()
	}()

	go func() {
		logger.Info("start shutdown gRPC server")
		grpcServer.Shutdown()
		logger.Info("end shutdown gRPC server")

		shutdownWG.Done()
	}()

	shutdownWG.Wait()

	if err := wg.Wait(); err != nil {
		logger.Error(err, "failed to run server")
		return exitNG
	}

	return exitOK
}
