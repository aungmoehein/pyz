package api

import (
	"net"

	"net/http"
)

// StatusResponseWriter wraps the http status code for logging
type StatusResponseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader wraps the http status code in ResponseWriter
func (w *StatusResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.status = code
}

// LanguageHeader used to pass language in request header
var LanguageHeader = "X-language"

// LoggingMiddleware log incoming http requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ws := &StatusResponseWriter{w, 200}
		next.ServeHTTP(ws, req)

		// write request log
		ip, _, _ := net.SplitHostPort(req.RemoteAddr)
		logger.Infof("%d %s %s %s", ws.status, ip, req.Method, req.URL)
	})
}

// PanicMiddleware write error response on panic errors
func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			language := getLocale(req)
			if r := recover(); r != nil {
				switch err := r.(type) {
				case Error:
					WriteError(w, err, language)

				default:
					WriteError(w, ErrUnexpected.Wraps(err), language)
				}
			}
		}()

		// serve next request
		next.ServeHTTP(w, req)
	})
}

// getLocale gets the relevant locale langauge from incoming request
func getLocale(req *http.Request) string {
	var language string
	// if header empty, get from url query
	if language = req.Header.Get(LanguageHeader); language == "" {
		language = req.URL.Query().Get("language")
	}

	// if still empty, use default language
	if language == "" {
		language = environ.DefaultLanguage
	}

	return language
}
