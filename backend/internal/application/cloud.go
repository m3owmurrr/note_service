package application

import (
	"cloud_technologies/internal/handlers"
	S3Storage "cloud_technologies/internal/storage/S3_storage"
	"log"
	"net/http"
)

const (
	addr = "0.0.0.0:8080"
	//addr = "localhost:8080"
)

func Run() {
	storage := S3Storage.NewS3Storage()
	noteHahdler := handlers.NewNotesHandler(storage)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.ChechHealth)
	mux.HandleFunc("POST /notes", noteHahdler.PostNoteHandler)
	mux.HandleFunc("GET /notes/{id}", noteHahdler.GetNoteHandler)

	handler := CORS(mux)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	log.Printf("server is running on %v...", addr)

	if err := server.ListenAndServe(); err != nil {
		log.Panicf("failed to start: %v", err)
	}
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("CORS: %s %s", r.Method, r.URL.Path)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
