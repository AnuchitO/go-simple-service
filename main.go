package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func l(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %s", r.RemoteAddr, r.Method, r.URL.Path, time.Since(start))
	}
}

func main() {
	host, err := os.Hostname()
	if err != nil {
		log.Println("error get hostname:", err)
	}
	http.HandleFunc("/", l(func(w http.ResponseWriter, r *http.Request) {
		w.Write(status("I am ok", host))
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

	log.Println("start server at port ", port)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func status(s string, h string) []byte {
	return []byte(`{"status": "` + s + `", "server": "` + h + `"}`)
}
