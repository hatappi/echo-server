package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	exitOK = iota
	exitNG
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	port := "3000"
	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}

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

	log.Printf("start http server: port is %s", port)
	http.ListenAndServe(":"+port, r)
	return exitOK
}

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Echo Server</title>
		<link rel="icon" href="data:image/x-icon;,">
	</head>
	<body>
		<h2>RequestPath</h2>
		{{.RequestPath}}

		<h2>Message</h2>
		{{.Message}}

		{{if ne (len .Meta) 0}}
			<h2>Metadata</h2>
			{{range $k, $v := .Meta}}
				<li>{{$k}}={{$v}}</li>
			{{end}}
		{{end}}

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

func handler(defaultMessage string, meta map[string]string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := defaultMessage
		if m := r.URL.Query().Get("message"); m != "" {
			msg = m
		}

		t, err := template.New("response").Parse(tpl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Message       string
			RequestPath   string
			RequestHeader http.Header
			Meta          map[string]string
		}{
			Message:       msg,
			RequestPath:   r.URL.Path,
			RequestHeader: r.Header,
			Meta:          meta,
		}

		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
