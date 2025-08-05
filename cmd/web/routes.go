package main

import (
	"net/http"
)

func (srv *server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", srv.makeHandler(srv.handleHome()))
	mux.HandleFunc("GET /page1", srv.makeHandler(srv.handlePage1()))
	mux.HandleFunc("GET /posts/view/{id}", srv.makeHandler(srv.postView()))

	return mux
}
