package app

import (
	"net/http"
	"os"
	"sync"
)

// rootIndexHandler serves the static landing page (GET/HEAD /).
// Resolution order: INDEX_HTML_PATH, then index.html (local dev), then /app/index.html (Docker).
func rootIndexHandler() http.HandlerFunc {
	var once sync.Once
	var body []byte
	var ok bool

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		once.Do(func() {
			paths := make([]string, 0, 3)
			if p := os.Getenv("INDEX_HTML_PATH"); p != "" {
				paths = append(paths, p)
			}
			paths = append(paths, "index.html", "/app/index.html")
			for _, p := range paths {
				b, err := os.ReadFile(p)
				if err != nil {
					continue
				}
				body = b
				ok = true
				return
			}
		})

		if !ok {
			http.Error(w, "landing page not found", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.Method == http.MethodHead {
			return
		}
		_, _ = w.Write(body)
	}
}
