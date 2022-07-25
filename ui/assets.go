package ui

import (
	"embed"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Assets hold the alice-lg frontend build
//go:embed build/*
var Assets embed.FS

// AssetsHTTPHandler handles HTTP request at a specific prefix.
// The prefix is usually /static.
func AssetsHTTPHandler(prefix string) http.Handler {
	handler := http.FileServer(http.FS(Assets))

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		reqPath := req.URL.Path
		rawPath := req.URL.RawPath

		if !strings.HasPrefix(reqPath, prefix) {
			handler.ServeHTTP(res, req)
			return
		}

		// Rewrite path
		reqPath = path.Join("build/", reqPath)
		rawPath = path.Join("build/", rawPath)

		// This is pretty much like the StripPrefix middleware,
		// from net/http, however we replace the prefix with `build/`.
		req1 := new(http.Request)
		*req1 = *req // clone request
		req1.URL = new(url.URL)
		*req1.URL = *req.URL

		req1.URL.Path = reqPath
		req1.URL.RawPath = rawPath

		handler.ServeHTTP(res, req1)
	})
}
