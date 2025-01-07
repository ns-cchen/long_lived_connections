package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a random delay between 1 and 3 seconds
	delay := rand.Intn(3) + 1
	podID := os.Getenv("POD_ID")
	if podID == "" {
		podID = "unknown" // Default to "unknown" if POD_ID is not set
	}

	// Log the delay and pod id
	fmt.Printf("Pod %s will delay the response for %d seconds.\n", podID, delay)

	// Introduce the delay
	time.Sleep(time.Duration(delay) * time.Second)

	// Write the response
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "Response from Pod %s, delayed by %d seconds", podID, delay)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func closeClientConnection(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		deadline, ok := r.Context().Value("key").(time.Time)
		if ok && time.Now().After(deadline) {
			w.Header().Set("Connection", "Close")
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func main() {
	// Get POD_ID from environment variable, default to "default" if not set
	podID := os.Getenv("POD_ID")
	if podID == "" {
		podID = "default"
	}

	// Set up the route for /api
	router := chi.NewRouter()
	router.Use(closeClientConnection)
	router.Get("/api", apiHandler)

	// Start the server
	fmt.Printf("Starting server on 0.0.0.0:5000 (POD_ID: %s)\n", podID)
	server := &http.Server{
		Addr: ":5000",
		Handler: router,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return context.WithValue(ctx, "key", time.Now().Add(120 * time.Second))
		},
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
