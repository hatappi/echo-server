package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-logr/logr"

	pb "github.com/hatappi/echo-server/pb"
)

func sayHelloHandler(metadata map[string]string, logger logr.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.SayHelloRequest

		switch r.Method {
		case http.MethodGet:
		case http.MethodHead:
			q := r.URL.Query()
			req.Name = q.Get("name")
		case http.MethodPost:
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, fmt.Sprintf("%s is not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}

		reqMeta := make(map[string]string)
		for k, vals := range r.Header {
			reqMeta[k] = strings.Join(vals, ",")
		}

		logger.Info("sayHello request received", "name", req.GetName(), "protocol", r.Proto, "path", r.URL.Path, "method", r.Method)

		resp := &pb.SayHelloResponse{
			Body:             "Hello " + req.GetName(),
			RequestMetadata:  reqMeta,
			ResponseMetadata: metadata,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type httpServer struct {
	logger logr.Logger

	server *http.Server
}

func NewHTTPServer(metadata map[string]string, logger logr.Logger) *httpServer {
	r := chi.NewRouter()
	r.HandleFunc("/*", sayHelloHandler(metadata, logger))

	server := &http.Server{
		Handler: r,
	}

	return &httpServer{
		logger: logger,

		server: server,
	}
}

func (hs *httpServer) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	hs.logger.Info("start HTTP server", "addr", addr)

	if err := hs.server.Serve(lis); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (hs *httpServer) Shutdown(ctx context.Context) error {
	if err := hs.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
