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

func (app *application) makeHandler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func (app *application) handleHomeGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.newTemplateData(r)
		data.ActivePage = "home"

		app.render(w, r, http.StatusOK, "home.tmpl", data)
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

func (app *application) handlePostCreateGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.newTemplateData(r)
		data.ActivePage = "create"
		data.Form = postCreateForm{}

		app.render(w, r, http.StatusOK, "create.tmpl", data)
	}
}

func (app *application) handlePostCreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 2 MiB
		r.Body = http.MaxBytesReader(w, r.Body, 2<<20)

		var form postCreateForm
		err := app.decodePostForm(r, &form)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
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
			data := app.newTemplateData(r)
			data.Form = form
			data.ActivePage = "create"
			app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
			return
		}

		s, err := app.posts.Insert(form.Title, form.Content, form.Excerpt, form.Author, form.Slug)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		app.sessionMGR.Put(r.Context(), "flash", "Post successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/posts/%s", s), http.StatusSeeOther)
	}
}

func (app *application) handlePostUpdateGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		post, err := app.posts.Get(slug)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.NotFound(w, r)
			} else {
				app.serverError(w, r, err)
			}
			return
		}

		data := app.newTemplateData(r)
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

		app.render(w, r, http.StatusOK, "update.tmpl", data)
	}
}

func (app *application) handlePostUpdatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		// 2 MiB
		r.Body = http.MaxBytesReader(w, r.Body, 2<<20)

		var form postCreateForm
		err = app.decodePostForm(r, &form)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
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
			data := app.newTemplateData(r)
			data.Form = form
			data.ActivePage = "post"
			app.render(w, r, http.StatusUnprocessableEntity, "update.tmpl", data)
			return
		}

		s, err := app.posts.Update(id, form.Title, form.Content, form.Excerpt, form.Author, form.Slug)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		app.sessionMGR.Put(r.Context(), "flash", "Post successfully updated!")

		http.Redirect(w, r, fmt.Sprintf("/posts/%s", s), http.StatusSeeOther)
	}
}

func (app *application) handlePostDeletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		err = app.posts.Delete(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		app.sessionMGR.Put(r.Context(), "flash", "Post successfully deleted!")

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}

func (app *application) handlePostsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := app.posts.Latest()
		if err != nil {
			app.serverError(w, r, err)
		}

		data := app.newTemplateData(r)
		data.ActivePage = "posts"
		data.Posts = posts

		app.render(w, r, http.StatusOK, "posts.tmpl", data)
	}
}

func (app *application) handlePostViewGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		post, err := app.posts.Get(slug)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.NotFound(w, r)
			} else {
				app.serverError(w, r, err)
			}
			return
		}

		htmlContent, err := app.renderMarkdownToHTML(post.Content)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		data := app.newTemplateData(r)
		data.ActivePage = "posts"
		data.Post = post
		data.HTMLContent = template.HTML(htmlContent)

		app.render(w, r, http.StatusOK, "post.tmpl", data)
	}
}
