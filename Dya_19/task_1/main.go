package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	quary := r.URL.Query()
	q := quary.Get("q")
	limit := quary.Get("limit")

	if q == "" {
		http.Error(w, "Parametr q is missing", http.StatusBadRequest)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		limitInt = 10
	}

	fmt.Fprintf(w, "Посик по запросу: %s, лимит %d", q, limitInt)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	quary := r.URL.Query()
	limit := quary.Get("limit")

	limitInt, _ := strconv.Atoi(limit)
	fmt.Fprintf(w, "Ебать тебя конопатый, хули ты здесь забыл еблан %d", limitInt)
}

func main() {

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/search", searchHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		fmt.Println("Server running on http://localhost:8080")
		fmt.Println("Example: http://localhost:8080/search?q=golang&limit=10")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	fmt.Println("\nPrint Enter to shutdown...")
	fmt.Scanln()

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Shutdown error: %v\n", err)
	}

	fmt.Println("Server successfully stoped")
}
