package server

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/roma-glushko/bp/internal/server/handlers"
	"github.com/roma-glushko/bp/internal/storage"
)

func NewMux(frontendFS embed.FS, store storage.Store) *http.ServeMux {
	mux := http.NewServeMux()

	sessions := &handlers.SessionHandler{Store: store}
	settings := &handlers.SettingsHandler{Store: store}
	reports := &handlers.ReportHandler{Store: store}

	mux.HandleFunc("GET /api/sessions", sessions.List)
	mux.HandleFunc("POST /api/sessions", sessions.Create)
	mux.HandleFunc("GET /api/sessions/{id}", sessions.Get)
	mux.HandleFunc("PUT /api/sessions/{id}", sessions.Update)
	mux.HandleFunc("DELETE /api/sessions/{id}", sessions.Delete)

	mux.HandleFunc("GET /api/settings", settings.Get)
	mux.HandleFunc("PUT /api/settings", settings.Update)

	mux.HandleFunc("GET /api/reports/preview", reports.Preview)

	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		slog.Error("failed to create sub filesystem", "error", err)
		return mux
	}

	fileServer := http.FileServer(http.FS(distFS))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			fileServer.ServeHTTP(w, r)
			return
		}

		f, err := distFS.Open(path[1:])
		if err != nil {
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)

			return
		}

		_ = f.Close()
		fileServer.ServeHTTP(w, r)
	})

	return mux
}
