package reqinspect

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type RequestInfo struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func Run() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read at most 10KB of the request body
		body := make([]byte, 10*1024)
		n, _ := io.ReadFull(r.Body, body) // ignore error

		reqInfo := RequestInfo{
			Method:  r.Method,
			Path:    r.URL.Path,
			Headers: r.Header,
			Body:    string(body[:n]),
		}

		slog.Info("Request received", slog.Any("request", reqInfo))
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		slog.Error(fmt.Sprintf("Failed to start server: %s", err.Error()))
	}
}
