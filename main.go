package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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

	httpPort := "3000"
	if p, ok := os.LookupEnv("ECHO_SERVER_HTTP_PORT"); ok {
		httpPort = p
	}
	wg.Go(func() error { return runHTTPServer(fmt.Sprintf("%s:%s", host, httpPort), meta, logger) })

	grpcPort := "5000"
	if p, ok := os.LookupEnv("ECHO_SERVER_GRPC_PORT"); ok {
		grpcPort = p
	}
	wg.Go(func() error { return runGRPCServer(fmt.Sprintf("%s:%s", host, grpcPort), meta, logger) })

	if err := wg.Wait(); err != nil {
		logger.Error(err, "failed to run server")
		return exitNG
	}

	return exitOK
}
