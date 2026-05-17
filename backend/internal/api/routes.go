package api

import "net/http"

// Register mounts API routes on the given mux.
func Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/config", ConfigHandler)
}
