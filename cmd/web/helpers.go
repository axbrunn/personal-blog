package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"

	"github.com/go-playground/form/v4"
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
	ts, ok := srv.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		srv.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		srv.serverError(w, r, err)
	}
}

func (srv *server) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),

		Flash: srv.sessionMGR.PopString(r.Context(), "flash"),
	}
}

func (srv *server) renderMarkdownToHTML(input string) (string, error) {
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	if err := md.Convert([]byte(input), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (srv *server) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request, in the same way that we did in our
	// snippetCreatePost handler.
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter.
	err = srv.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		// For all other errors, return them as normal.
		return err
	}
	return nil
}
