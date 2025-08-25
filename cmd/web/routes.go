package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionMGR.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.makeHandler(app.handleHomeGet())))
	mux.Handle("GET /posts", dynamic.ThenFunc(app.makeHandler(app.handlePostsGet())))
	mux.Handle("GET /posts/{slug}", dynamic.ThenFunc(app.makeHandler(app.handlePostViewGet())))
	mux.Handle("GET /posts/create", dynamic.ThenFunc(app.makeHandler(app.handlePostCreateGet())))
	mux.Handle("POST /posts/create", dynamic.ThenFunc(app.makeHandler(app.handlePostCreatePost())))
	mux.Handle("POST /posts/delete/{id}", dynamic.ThenFunc(app.makeHandler(app.handlePostDeletePost())))
	mux.Handle("GET /posts/update/{slug}", dynamic.ThenFunc(app.makeHandler(app.handlePostUpdateGet())))
	mux.Handle("POST /posts/update/{id}", dynamic.ThenFunc(app.makeHandler(app.handlePostUpdatePost())))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
