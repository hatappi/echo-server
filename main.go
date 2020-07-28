package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handler)

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
	</head>
	<body>
		<h2>Message</h2>
		{{.Message}}
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

func handler(w http.ResponseWriter, r *http.Request) {
	msg := "hello"
	if m, ok := os.LookupEnv("MESSAGE"); ok {
		msg = m
	}

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
		RequestHeader http.Header
	}{
		Message:       msg,
		RequestHeader: r.Header,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
