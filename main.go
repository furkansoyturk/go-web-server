package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("/logo", http.FileServer(http.Dir("/assets")))
	mux.HandleFunc("/app", app)
	mux.HandleFunc("/app/assets", assets)
	mux.HandleFunc("/healthz", healthzHandler)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port: %s \n", port)
	log.Fatal(server.ListenAndServe())

}

func assets(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, req, "./assets/index.html")
}
func app(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, req, "index.html")
}
func healthzHandler(w http.ResponseWriter, req *http.Request) {
	ok := []byte("OK")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(ok)
}
