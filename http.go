package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func runHTTPServer(addr string) error {
	msg := "hello"
	if m, ok := os.LookupEnv("MESSAGE"); ok {
		msg = m
	}

	meta := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], "ECHO_SERVER_META_") {
			meta[strings.TrimPrefix(pair[0], "ECHO_SERVER_META_")] = pair[1]
		}
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/*", handler(msg, meta))

	log.Printf("start http server: addr is %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		return err
	}

	return nil
}

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Echo Server</title>
		<link rel="icon" href="data:image/x-icon;,">
		{{if ne (len .Style) 0}}
			<style type="text/css">
			{{.Style}}
			</style>
		{{end}}
	</head>
	<body>
		<h2>RequestPath</h2>
		{{.RequestPath}}

		<h2>ClientIP</h2>
		{{.ClientIP}}

		<h2>Message</h2>
		{{.Message}}

		{{if ne (len .Meta) 0}}
			<h2>Metadata</h2>
			{{range $k, $v := .Meta}}
				<li>{{$k}}={{$v}}</li>
			{{end}}
		{{end}}

		<h2>Host</h2>
		<div>{{.Host}}</div>

		<h2>Request Header</h2>
		{{range $k, $arr := .RequestHeader}}
			<div><b>{{$k}}</b></div>
			<ul>
				{{range $arr}}
					<li>{{.}}</li>
				{{end}}
			</ul>
		{{end}}
	</body>
</html>`

type responseData struct {
	Message       string            `json:"message"`
	RequestPath   string            `json:"request_path"`
	RequestHeader http.Header       `json:"request_header"`
	Host          string            `json:"host"`
	Meta          map[string]string `json:"meta"`
	ClientIP      string            `json:"client_ip"`
}

func handler(defaultMessage string, meta map[string]string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := defaultMessage
		if m := r.URL.Query().Get("message"); m != "" {
			msg = m
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := &responseData{
			Message:       msg,
			RequestPath:   r.URL.Path,
			RequestHeader: r.Header,
			Host:          r.Host,
			Meta:          meta,
			ClientIP:      ip,
		}

		if r.Header.Get("Content-Type") == "application/json" {
			b, err := json.Marshal(data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(b)
		} else {
			t, err := template.New("response").Parse(tpl)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			templateData := struct {
				Message       string
				RequestPath   string
				Host          string
				RequestHeader http.Header
				Meta          map[string]string
				Style         template.CSS
				ClientIP      string
			}{
				Message:       data.Message,
				RequestPath:   data.RequestPath,
				RequestHeader: data.RequestHeader,
				Host:          data.Host,
				Meta:          data.Meta,
				Style:         template.CSS(os.Getenv("STYLE")),
				ClientIP:      data.ClientIP,
			}

			err = t.Execute(w, templateData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
