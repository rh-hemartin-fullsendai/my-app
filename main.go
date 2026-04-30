package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	counter   int
	counterMu sync.Mutex
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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
	counterMu.Lock()
	defer counterMu.Unlock()
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%d", counter)
}

func counterIncrementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	counterMu.Lock()
	defer counterMu.Unlock()
	counter++
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%d", counter)
}
