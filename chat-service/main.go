// package main

// import (
// 	"chat-service/db"
// 	"chat-service/handlers"
// 	"chat-service/middleware"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// func main() {
// 	// Load .env (for local dev)
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Println("No .env file found, relying on system environment variables.")
// 	}

// 	// Initialize DB
// 	err = db.Init()
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}

// 	// Routes
// 	mux := http.NewServeMux()

// 	// Ping test with JWT auth
// 	mux.Handle("/ping", middleware.Protected(http.HandlerFunc(handlers.PingWithAuth)))

// 	// Check if user follows back
// 	mux.Handle("/mutual-follow", middleware.Protected(http.HandlerFunc(handlers.CheckMutualFollow)))

// 	// WebSocket chat handler
// 	mux.Handle("/ws", middleware.Protected(http.HandlerFunc(handlers.HandleWebSocket)))

//		// Start server
//		port := os.Getenv("PORT")
//		if port == "" {
//			port = "8080"
//		}
//		fmt.Printf("Server running on port %s\n", port)
//		log.Fatal(http.ListenAndServe(":"+port, mux))
//	}
package main

import (
	"chat-service/db"
	"chat-service/handlers"
	"chat-service/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env (for local dev)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables.")
	}
	// Initialize database connection
	if err := db.Init(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	log.Println("Connected to PostgreSQL")

	// Create router
	r := mux.NewRouter()

	// // Public WebSocket route (auth handled inside handler)
	// r.HandleFunc("/ws", handlers.HandleWebSocket)

	r.Handle("/ws", middleware.Protected(http.HandlerFunc(handlers.HandleWebSocket)))

	// Protected routes
	r.Handle("/ping", middleware.Protected(http.HandlerFunc(handlers.PingWithAuth))).Methods("GET")
	r.Handle("/check-mutual-follow", middleware.Protected(http.HandlerFunc(handlers.CheckMutualFollow))).Methods("POST")
	r.Handle("/follow", middleware.Protected(http.HandlerFunc(handlers.FollowUser))).Methods("POST") // Assuming this is implemented

	// Start server
	log.Println("Chat service running on :8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}
