package server

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
)

func NewMux(frontendFS embed.FS) *http.ServeMux {
	mux := http.NewServeMux()

	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		slog.Error("failed to create sub filesystem", "error", err)
		return mux
	}

	fileServer := http.FileServer(http.FS(distFS))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly; fall back to index.html for SPA
		path := r.URL.Path
		if path == "/" {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Check if file exists in the embedded FS
		f, err := distFS.Open(path[1:]) // strip leading /
		if err != nil {
			// File not found — serve index.html for SPA routing
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}
		f.Close()
		fileServer.ServeHTTP(w, r)
	})

	return mux
}
