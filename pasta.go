package main

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var pasta atomic.Value

func init() {
	pasta.Store("")
}

//go:embed views/index.html
var indexHtml string
var indexTemplate = template.Must(template.New("index").Parse(indexHtml))

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		p := pasta.Load().(string)
		err := indexTemplate.Execute(w, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" && r.URL.Path == "/update" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p := r.PostFormValue("pasta")
		pasta.Store(p)
		http.Redirect(w, r, "/#saved", 302)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	srv := &http.Server{
		Addr:         ":3000",
		Handler:      http.HandlerFunc(handler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
