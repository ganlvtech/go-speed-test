package main

import (
	_ "embed"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

//go:embed index.html
var indexHTML []byte

var buf [1024 * 1024]byte

func main() {
	http.HandleFunc("/file.bin", func(w http.ResponseWriter, r *http.Request) {
		sizeStr := r.URL.Query().Get("size")
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			slog.Error("size query parameter error", "err", err, "sizeStr", sizeStr)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				slog.Error("Write error", "err", err)
				return
			}
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.Itoa(size))
		w.WriteHeader(http.StatusOK)
		for i := 0; i < size; i += len(buf) {
			remainingBytes := size - i
			if remainingBytes >= len(buf) {
				_, err := w.Write(buf[:])
				if err != nil {
					slog.Error("Write error", "err", err)
					return
				}
			} else {
				_, err := w.Write(buf[:remainingBytes])
				if err != nil {
					slog.Error("Write error", "err", err)
					return
				}
			}
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(len(indexHTML)))
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(indexHTML)
		if err != nil {
			slog.Error("Write error", "err", err)
			return
		}
	})
	listenAddr := ":8000"
	if len(os.Args) >= 2 {
		listenAddr = os.Args[1]
	}
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		slog.Error("ListenAndServe exit", "err", err)
		return
	}
}
