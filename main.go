package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
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

func toMB(b uint64) float64 {
	return float64(b) / float64(MB)
}

func main() {
	host, err := os.Hostname()
	if err != nil {
		log.Println("error get hostname:", err)
	}
	http.HandleFunc("/", l(func(w http.ResponseWriter, r *http.Request) {
		w.Write(status("I am ok", host))
	}))

	http.HandleFunc("/metrics", l(func(w http.ResponseWriter, r *http.Request) {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		alloc := toMB(mem.Alloc)
		totalAlloc := toMB(mem.TotalAlloc)
		sysAlloc := toMB(mem.Sys)
		heapInuse := toMB(mem.HeapInuse)
		heapIdle := toMB(mem.HeapIdle)
		heapReleased := toMB(mem.HeapReleased)
		stackInuse := toMB(mem.StackInuse)
		stackSys := toMB(mem.StackSys)

		w.Header().Set("Content-Type", "application/json")
		metrics := fmt.Sprintf(`{"size": "MB", "Alloc": %.2f, "TotalAlloc": %.2f, "Sys": %.2f, "HeapInuse": %.2f, "HeapIdle": %.2f, "HeapReleased": %.2f, "StackInuse": %.2f, "StackSys": %.2f}`, alloc, totalAlloc, sysAlloc, heapInuse, heapIdle, heapReleased, stackInuse, stackSys)
		w.Write([]byte(metrics))
	}))

	http.HandleFunc("/healthz", l(func(w http.ResponseWriter, r *http.Request) {
		w.Write(status("ok", host))
	}))

	http.HandleFunc("/liveness", l(func(w http.ResponseWriter, r *http.Request) {
		w.Write(status("live", host))
	}))

	http.HandleFunc("/readiness", l(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.Write(status("ready", host))
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("start server at host %s port %s\n", host, port)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func status(s string, h string) []byte {
	return []byte(`{"status": "` + s + `", "server": "` + h + `"}`)
}
