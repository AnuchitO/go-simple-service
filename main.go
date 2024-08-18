package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func l(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %s", r.RemoteAddr, r.Method, r.URL.Path, time.Since(start))
	}
}

type Size uint64

const (
	Byte Size = 1 << (10 * iota)
	KB
	MB
)

func megabytes(b uint64) float64 {
	return float64(b) / float64(MB)
}

func toMB(b uint64) string {
	return fmt.Sprintf("%.2f MB", megabytes(b))
}

// go build -ldflags "-X main.version=v1.0.0 -X main.commit=123456"
var commit string
var version string = "v0.0.0"

var first = true

func main() {
	host, err := os.Hostname()
	if err != nil {
		log.Println("error get hostname:", err)
	}

	http.HandleFunc("/", l(func(w http.ResponseWriter, r *http.Request) {
		w.Write(status("I am ok", host))
	}))

	http.HandleFunc("/metrics", l(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(metrics()))
	}))

	http.HandleFunc("/versions", l(func(w http.ResponseWriter, r *http.Request) {
		v := versions()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(v))
	}))

	http.HandleFunc("/liveness", l(func(w http.ResponseWriter, r *http.Request) {
		w.Write(status("live", host))
	}))

	http.HandleFunc("/readiness", l(func(w http.ResponseWriter, r *http.Request) {
		if first {
			time.Sleep(5 * time.Second)
			first = false
		}

		w.Write(status("ready", host))
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// TODO: graceful shutdown
	// go gracefully(srv)

	log.Printf("start server at host %s port %s\n", host, port)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func status(s string, h string) []byte {
	return []byte(`{"status": "` + s + `", "server": "` + h + `"}`)
}

func metrics() string {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return fmt.Sprintf(`{
  "alloc": %q,
  "totalAlloc": %q,
  "sysAlloc": %q,
  "heapInuse": %q,
  "heapIdle": %q,
  "heapReleased": %q,
  "stackInuse": %q,
  "stackSys": %q
}`, toMB(mem.Alloc),
		toMB(mem.TotalAlloc),
		toMB(mem.Sys),
		toMB(mem.HeapInuse),
		toMB(mem.HeapIdle),
		toMB(mem.HeapReleased),
		toMB(mem.StackInuse),
		toMB(mem.StackSys),
	)
}

func versions() string {
	goVersion := runtime.Version()
	osVersion := runtime.GOOS
	arch := runtime.GOARCH
	host, _ := os.Hostname()

	return fmt.Sprintf(`{"go": "%s", "os": "%s", "arch": "%s", "host": "%s", "commit": "%s", "version": "%s"}`, goVersion, osVersion, arch, host, commit, version)
}

func gracefully(srv *http.Server) {
	{
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Kill, os.Interrupt)
		defer cancel()
		<-ctx.Done()
	}

	d := time.Duration(5 * time.Second)
	slog.Info(fmt.Sprintf("Server shutting down within %d ...\n", d))

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Info("shutting down error:" + err.Error())
	}
}
