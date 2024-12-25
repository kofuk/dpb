package main

import (
	"net/http"
	"os"
	"strconv"
)

func main() {
	statusCode := 200
	if os.Getenv("STATUS_CODE") != "" {
		var err error
		statusCode, err = strconv.Atoi(os.Getenv("STATUS_CODE"))
		if err != nil {
			panic(err)
		}
	}

	body := os.Getenv("MESSAGE_BODY")

	contentType := os.Getenv("CONTENT_TYPE")
	if contentType == "" {
		contentType = "text/plain"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	})

	http.ListenAndServe(":8000", nil)
}
