package server

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/api"
	"github.com/gabriel/ttrpg-toolkit/backend/ui"
)

// New returns an HTTP handler for the application.
func New() http.Handler {
	mux := http.NewServeMux()
	api.Register(mux)

	staticFS, err := fs.Sub(ui.Static, "static")
	if err != nil {
		panic("embedded static files: " + err.Error())
	}

	fileServer := http.FileServer(http.FS(staticFS))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			http.ServeFileFS(w, r, staticFS, "index.html")
			return
		}

		clean := path
		if strings.Contains(clean, "..") {
			http.NotFound(w, r)
			return
		}

		if f, err := staticFS.Open(clean); err == nil {
			f.Close()
			r.URL.Path = "/" + clean
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFileFS(w, r, staticFS, "index.html")
	})

	return mux
}
