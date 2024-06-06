package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type Request struct {
	Body string `json:"body"`
}

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
	mux.HandleFunc("POST /api/validate_chirp", validateLength)
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
func validateLength(w http.ResponseWriter, r *http.Request) {
	var req Request
	request, err := io.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong"})
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(request, &req)
	req = censorRequestBody(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong"})
		return
	}
	if len(req.Body) <= 140 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"cleaned_body": req.Body})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Chirp is too long"})
	}

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

func censorRequestBody(req Request) Request {
	words := strings.Split(req.Body, " ")
	cleanedReq := []string{}
	for _, s := range words {
		switch strings.ToLower(s) {
		case "kerfuffle":
			s = "****"
		case "sharbert":
			s = "****"
		case "fornax":
			s = "****"
		}
		cleanedReq = append(cleanedReq, s)
	}
	req.Body = strings.Join(cleanedReq, " ")
	return req
}
