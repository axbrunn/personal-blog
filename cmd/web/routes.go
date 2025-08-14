package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (srv *server) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", srv.makeHandler(srv.handleHomeGet()))
	mux.HandleFunc("GET /posts", srv.makeHandler(srv.handlePostsGet()))
	mux.HandleFunc("GET /posts/{slug}", srv.makeHandler(srv.handlePostViewGet()))
	mux.HandleFunc("GET /posts/create", srv.makeHandler(srv.handlePostCreateGet()))
	mux.HandleFunc("POST /posts/create", srv.makeHandler(srv.handlePostCreatePost()))
	mux.HandleFunc("POST /posts/delete/{id}", srv.makeHandler(srv.handlePostDeletePost()))

	standard := alice.New(srv.recoverPanic, srv.logRequest, commonHeaders)

	return standard.Then(mux)
}
