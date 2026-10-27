package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	dbPath := "enquiries.db"
	db, err := openDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	clientDir := os.Getenv("CLIENT_DIR")
	if clientDir == "" {
		clientDir = "../client/dist"
	}
	if _, err := os.Stat(filepath.Join(clientDir, "index.html")); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	client := http.FileServer(http.Dir(clientDir))
	mux.Handle("GET /assets/", client)
	serveClient := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(clientDir, "index.html"))
	}
	mux.HandleFunc("GET /{$}", serveClient)
	mux.HandleFunc("GET /staff", serveClient)
	registerRoutes(mux, db)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
