package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("invoke: received a request")

	path := r.URL.Path

	input := strings.Trim(path, "/")

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/app/hw", input)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		writeResponse(w, 500, err.Error())
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		writeResponse(w, 500, err.Error())
		return
	}

	if err := cmd.Start(); err != nil {
		writeResponse(w, 500, err.Error())
		return
	}

	// fmt.Fprintf(stdin, "%s\n", input) // io.WriteString(stdin, input)

	seout, _ := ioutil.ReadAll(stderr)
	stout, _ := ioutil.ReadAll(stdout)

	if err := cmd.Wait(); err != nil {
		writeResponse(w, 500, string(seout))
		return
	}

	if len(seout) > 0 {
		writeResponse(w, 400, string(seout))
		return
	}
	writeResponse(w, 200, string(stout))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler)
	logHandler := loggingHandler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("invoker: listening on %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), logHandler))
}

type statusRecorder struct {
	http.ResponseWriter
	status    int
	byteCount int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *statusRecorder) Write(p []byte) (int, error) {
	bc, err := rec.ResponseWriter.Write(p)
	rec.byteCount += bc

	return bc, err
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rec := statusRecorder{w, 200, 0}
		next.ServeHTTP(&rec, req)
		remoteAddr := req.Header.Get("X-Forwarded-For")
		if remoteAddr == "" {
			remoteAddr = req.RemoteAddr
		}
		ua := req.Header.Get("User-Agent")

		log.Printf("%s - \"%s %s %s\" (%s) %d %d \"%s\"", remoteAddr, req.Method, req.URL.Path, req.Proto, req.Host, rec.status, rec.byteCount, ua)
	})
}

func writeResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(message))
}
