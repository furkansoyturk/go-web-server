package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	apiConfig := ApiConfig{
		FileServerHits: 0,
	}
	const filepathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("/api/reset", apiConfig.handerReset)
	mux.HandleFunc("/admin/metrics", apiConfig.adminMiddlewareMetricsInc)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (apiConfig *ApiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", apiConfig.FileServerHits)))
}

func (apiConfig *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiConfig.FileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (apiConfig *ApiConfig) adminMiddlewareMetricsInc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	t, _ := template.ParseFiles("admin_index.html")
	items := struct {
		Value int
	}{
		Value: apiConfig.FileServerHits,
	}
	t.Execute(w, items)
}
func (apiConfig *ApiConfig) handerReset(w http.ResponseWriter, r *http.Request) {
	apiConfig.FileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
