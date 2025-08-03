package main

import (
	"fmt"
	"net/http"
)

func (srv *server) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.RequestURI
	)

	srv.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (srv *server) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (srv *server) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Retrieve the appropriate template set from the cache based on the page
	// name (like 'home.tmpl'). If no entry exists in the cache with the
	// provided name, then create a new error and call the serverError() helper
	// method that we made earlier and return.
	ts, ok := srv.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		srv.serverError(w, r, err)
		return
	}
	// Write out the provided HTTP status code ('200 OK', '400 Bad Request' etc).
	w.WriteHeader(status)
	// Execute the template set and write the response body. Again, if there
	// is any error we call the serverError() helper.
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		srv.serverError(w, r, err)
	}
}
