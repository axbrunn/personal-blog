package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (srv *server) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(srv.sessionMGR.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(srv.makeHandler(srv.handleHomeGet())))
	mux.Handle("GET /posts", dynamic.ThenFunc(srv.makeHandler(srv.handlePostsGet())))
	mux.Handle("GET /posts/{slug}", dynamic.ThenFunc(srv.makeHandler(srv.handlePostViewGet())))
	mux.Handle("GET /posts/create", dynamic.ThenFunc(srv.makeHandler(srv.handlePostCreateGet())))
	mux.Handle("POST /posts/create", dynamic.ThenFunc(srv.makeHandler(srv.handlePostCreatePost())))
	mux.Handle("POST /posts/delete/{id}", dynamic.ThenFunc(srv.makeHandler(srv.handlePostDeletePost())))
	mux.Handle("GET /posts/update/{slug}", dynamic.ThenFunc(srv.makeHandler(srv.handlePostUpdateGet())))
	mux.Handle("POST /posts/update/{id}", dynamic.ThenFunc(srv.makeHandler(srv.handlePostUpdatePost())))

	standard := alice.New(srv.recoverPanic, srv.logRequest, commonHeaders)

	return standard.Then(mux)
}
