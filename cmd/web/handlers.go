package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/axbrunn/http_web/internals/models"
)

func (srv *server) makeHandler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func (srv *server) handleHomeGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := srv.newTemplateData(r)
		data.ActivePage = "home"

		srv.render(w, r, http.StatusOK, "home.tmpl", data)
	}
}

func (srv *server) handlePostCreateGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := srv.newTemplateData(r)
		data.ActivePage = "create"

		srv.render(w, r, http.StatusOK, "create.tmpl", data)
	}
}

func (srv *server) handlePostCreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			srv.clientError(w, http.StatusBadRequest)
			return
		}

		author := r.PostForm.Get("author")
		slug := r.PostForm.Get("slug")
		title := r.PostForm.Get("title")
		excerpt := r.PostForm.Get("excerpt")
		content := r.PostForm.Get("content")

		s, err := srv.posts.Insert(title, content, excerpt, author, slug)
		if err != nil {
			srv.serverError(w, r, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/posts/%s", s), http.StatusSeeOther)
	}
}

func (srv *server) handlePostUpdateGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := srv.newTemplateData(r)
		data.ActivePage = "update"

		srv.render(w, r, http.StatusOK, "update.tmpl", data)
	}
}

func (srv *server) handlePostUpdatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (srv *server) handlePostDeletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			srv.clientError(w, http.StatusBadRequest)
			return
		}

		err = srv.posts.Delete(id)
		if err != nil {
			srv.serverError(w, r, err)
			return
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}

func (srv *server) handlePostsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := srv.posts.Latest()
		if err != nil {
			srv.serverError(w, r, err)
		}

		data := srv.newTemplateData(r)
		data.ActivePage = "posts"
		data.Posts = posts

		srv.render(w, r, http.StatusOK, "posts.tmpl", data)
	}
}

func (srv *server) handlePostViewGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		post, err := srv.posts.Get(slug)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.NotFound(w, r)
			} else {
				srv.serverError(w, r, err)
			}
			return
		}

		htmlContent, err := srv.renderMarkdownToHTML(post.Content)
		if err != nil {
			srv.serverError(w, r, err)
			return
		}

		data := srv.newTemplateData(r)
		data.ActivePage = "posts"
		data.Post = post
		data.HTMLContent = template.HTML(htmlContent)

		srv.render(w, r, http.StatusOK, "post.tmpl", data)
	}
}
