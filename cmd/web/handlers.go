package main

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"

	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/axbrunn/http_web/internals/models"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/parser"
)

func (s *server) makeHandler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func (srv *server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "Go")

		data := srv.newTemplateData(r)
		data.ActivePage = "home"

		// Use the new render helper.
		srv.render(w, r, http.StatusOK, "home.tmpl", data)
	}
}

func (srv *server) handlePosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "Go")

		posts, err := srv.posts.Latest()
		if err != nil {
			srv.serverError(w, r, err)
		}

		data := srv.newTemplateData(r)
		data.ActivePage = "posts"
		data.Posts = posts

		// Use the new render helper.
		srv.render(w, r, http.StatusOK, "posts.tmpl", data)
	}
}

func (srv *server) postView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		// Use the SnippetModel's Get() method to retrieve the data for a
		// specific record based on its ID. If no matching record is found,
		// return a 404 Not Found response.
		post, err := srv.posts.Get(slug)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.NotFound(w, r)
			} else {
				srv.serverError(w, r, err)
			}
			return
		}

		// Markdown naar HTML
		var buf bytes.Buffer
		// Goldmark configureren
		md := goldmark.New(
			goldmark.WithExtensions(
				highlighting.NewHighlighting(
					highlighting.WithStyle("dracula"),
					highlighting.WithFormatOptions(
						chromahtml.WithLineNumbers(true),
					),
				),
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
		)

		if err := md.Convert([]byte(post.Content), &buf); err != nil {
			srv.serverError(w, r, err)
			return
		}

		data := srv.newTemplateData(r)
		data.ActivePage = "posts"
		data.Post = post
		data.HTMLContent = template.HTML(buf.String())

		srv.render(w, r, http.StatusOK, "post.tmpl", data)
	}
}
