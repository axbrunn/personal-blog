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

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.RequestURI
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	if _, err := buf.WriteTo(w); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),

		Flash: app.sessionMGR.PopString(r.Context(), "flash"),
	}
}

func (app *application) renderMarkdownToHTML(input string) (string, error) {
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

func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request, in the same way that we did in our
	// snippetCreatePost handler.
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter.
	err = app.formDecoder.Decode(dst, r.PostForm)
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
