package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/axbrunn/http_web/internals/models"
	"github.com/axbrunn/http_web/internals/validator"
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

type postCreateForm struct {
	ID                  int    `form:"-"`
	Author              string `form:"author"`
	Slug                string `form:"slug"`
	Title               string `form:"title"`
	Excerpt             string `form:"excerpt"`
	Content             string `form:"content"`
	validator.Validator `form:"-"`
}

func (srv *server) handlePostCreateGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := srv.newTemplateData(r)
		data.ActivePage = "create"
		data.Form = postCreateForm{}

		srv.render(w, r, http.StatusOK, "create.tmpl", data)
	}
}

func (srv *server) handlePostCreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 2 MiB
		r.Body = http.MaxBytesReader(w, r.Body, 2<<20)

		var form postCreateForm
		err := srv.decodePostForm(r, &form)
		if err != nil {
			srv.clientError(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Author), "author", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Author, 100), "author", "This field cannot be more than 100 characters long")
		form.CheckField(validator.NotBlank(form.Slug), "slug", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Slug, 100), "slug", "This field cannot be more than 100 characters long")
		form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Title, 255), "title", "This field cannot be more than 255 characters long")
		form.CheckField(validator.NotBlank(form.Excerpt), "excerpt", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Excerpt, 255), "excerpt", "This field cannot be more than 255 characters long")
		form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

		if !form.Valid() {
			data := srv.newTemplateData(r)
			data.Form = form
			data.ActivePage = "create"
			srv.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
			return
		}

		s, err := srv.posts.Insert(form.Title, form.Content, form.Excerpt, form.Author, form.Slug)
		if err != nil {
			srv.serverError(w, r, err)
			return
		}

		srv.sessionMGR.Put(r.Context(), "flash", "Post successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/posts/%s", s), http.StatusSeeOther)
	}
}

func (srv *server) handlePostUpdateGet() http.HandlerFunc {
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

		data := srv.newTemplateData(r)
		data.ActivePage = "posts"
		data.Post = post
		data.Form = postCreateForm{
			ID:      post.ID,
			Author:  post.Author,
			Slug:    post.Slug,
			Title:   post.Title,
			Excerpt: post.Excerpt,
			Content: post.Content,
		}

		srv.render(w, r, http.StatusOK, "update.tmpl", data)
	}
}

func (srv *server) handlePostUpdatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			srv.clientError(w, http.StatusBadRequest)
			return
		}

		// 2 MiB
		r.Body = http.MaxBytesReader(w, r.Body, 2<<20)

		var form postCreateForm
		err = srv.decodePostForm(r, &form)
		if err != nil {
			srv.clientError(w, http.StatusBadRequest)
			return
		}

		form.ID = id
		form.CheckField(validator.NotBlank(form.Author), "author", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Author, 100), "author", "This field cannot be more than 100 characters long")
		form.CheckField(validator.NotBlank(form.Slug), "slug", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Slug, 100), "slug", "This field cannot be more than 100 characters long")
		form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Title, 255), "title", "This field cannot be more than 255 characters long")
		form.CheckField(validator.NotBlank(form.Excerpt), "excerpt", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Excerpt, 255), "excerpt", "This field cannot be more than 255 characters long")
		form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

		if !form.Valid() {
			data := srv.newTemplateData(r)
			data.Form = form
			data.ActivePage = "create"
			srv.render(w, r, http.StatusUnprocessableEntity, "update.tmpl", data)
			return
		}

		s, err := srv.posts.Update(id, form.Title, form.Content, form.Excerpt, form.Author, form.Slug)
		if err != nil {
			srv.serverError(w, r, err)
			return
		}

		srv.sessionMGR.Put(r.Context(), "flash", "Post successfully updated!")

		http.Redirect(w, r, fmt.Sprintf("/posts/%s", s), http.StatusSeeOther)
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

		srv.sessionMGR.Put(r.Context(), "flash", "Post successfully deleted!")

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
