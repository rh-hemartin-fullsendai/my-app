package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "counter.db"
	}

	if err := initDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := closeDB(); err != nil {
			log.Printf("closing database: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/counter", counterHandler)
	mux.HandleFunc("/counter/increment", counterIncrementHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Listen for shutdown signals so deferred cleanup (closeDB) runs.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	fmt.Printf("Listening on :%s\n", port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	value, err := getCounter()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error getting counter: %v", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%d", value)
}

func counterIncrementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	value, err := incrementCounter()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error incrementing counter: %v", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%d", value)
}
