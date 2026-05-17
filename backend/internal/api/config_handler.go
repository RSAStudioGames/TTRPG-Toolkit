package api

import "net/http"

// ConfigData is exposed to the frontend at runtime.
type ConfigData struct {
	APIBaseURL string `json:"api_base_url"`
	WSURL      string `json:"ws_url"`
}

// ConfigHandler serves GET /api/config.
func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}
	WriteSuccess(w, ConfigData{
		APIBaseURL: "/api",
		WSURL:      "",
	})
}
