package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	defer closeDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/counter", counterHandler)
	mux.HandleFunc("/counter/increment", counterIncrementHandler)

	fmt.Printf("Listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
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
