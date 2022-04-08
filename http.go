package main

import (
	"encoding/json"
	"fmt"
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

		logger.Info("sayHello request received", "name", req.GetName(), "protocol", "HTTP")

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

func runHTTPServer(addr string, metadata map[string]string, logger logr.Logger) error {
	r := chi.NewRouter()
	r.HandleFunc("/*", sayHelloHandler(metadata, logger))

	logger.Info("start HTTP server", "addr", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		return err
	}

	return nil
}
