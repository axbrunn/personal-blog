package main

import (
	"net/http"
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
