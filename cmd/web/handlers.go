package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/axbrunn/http_web/internals/models"
)

func (s *server) makeHandler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func (srv *server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "Go")

		// Use the new render helper.
		srv.render(w, r, http.StatusOK, "home.tmpl", templateData{
			ActivePage: "home",
		})
	}
}

func (srv *server) handlePage1() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "Go")

		// Use the new render helper.
		srv.render(w, r, http.StatusOK, "page1.tmpl", templateData{
			ActivePage: "page1",
		})
	}
}

func (srv *server) postView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}
		// Use the SnippetModel's Get() method to retrieve the data for a
		// specific record based on its ID. If no matching record is found,
		// return a 404 Not Found response.
		post, err := srv.posts.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.NotFound(w, r)
			} else {
				srv.serverError(w, r, err)
			}
			return
		}
		// Write the snippet data as a plain-text HTTP response body.
		fmt.Fprintf(w, "%+v", post)
	}
}
