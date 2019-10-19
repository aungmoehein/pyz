package api

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"encoding/json"
	"net/http/httptest"

	"github.com/go-chi/chi"
)

// CallTestAPI is Helper function to test a handler with params json body and
// save resulting JSON output in out
// FIXME: PLEASE REFACTOR THIS
func CallTestAPI(router chi.Router, method, key, path string, params, out interface{}) {
	var err error
	var body io.Reader
	var data []byte

	if params != nil {
		if data, err = json.Marshal(params); err != nil {
			panic(err)
		}
		body = bytes.NewReader(data)
	}

	var recorder = httptest.NewRecorder()
	var request = httptest.NewRequest(method, path, body)

	request.Header.Set("Authorization", "Bearer "+key)
	request.Header.Set("InstanceID", environ.ID) // Mainly used to detect tests
	router.ServeHTTP(recorder, request)

	if out != nil {
		data = recorder.Body.Bytes()
		if err = json.Unmarshal(data, out); err != nil {
			panic(err)
		}
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
