// // main.go
// package main

// import (
// 	"chat-service/db"
// 	"chat-service/handlers"
// 	"chat-service/middleware"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/mux"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	// Load .env (for local dev)
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Println("No .env file found, relying on system environment variables.")
// 	}
// 	// Initialize database connection
// 	if err := db.Init(); err != nil {
// 		log.Fatal("Failed to connect to DB:", err)
// 	}
// 	log.Println("Connected to PostgreSQL")

// 	// Create router
// 	r := mux.NewRouter()

// 	// // Public WebSocket route (auth handled inside handler)
// 	// r.HandleFunc("/ws", handlers.HandleWebSocket)

// 	r.Handle("/ws", middleware.Protected(http.HandlerFunc(handlers.HandleWebSocket)))

// 	// Protected routes
// 	r.Handle("/ping", middleware.Protected(http.HandlerFunc(handlers.PingWithAuth))).Methods("GET")
// 	r.Handle("/check-mutual-follow", middleware.Protected(http.HandlerFunc(handlers.CheckMutualFollow))).Methods("POST")
// 	r.Handle("/follow", middleware.Protected(http.HandlerFunc(handlers.FollowUser))).Methods("POST") // Assuming this is implemented

// 	// Start server
// 	log.Println("Chat service running on :8083")
// 	log.Fatal(http.ListenAndServe(":8083", r))
// }

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

	// WebSocket route
	r.Handle("/ws", middleware.Protected(http.HandlerFunc(handlers.HandleWebSocket)))

	// Protected routes
	r.Handle("/ping", middleware.Protected(http.HandlerFunc(handlers.PingWithAuth))).Methods("GET")
	r.Handle("/check-mutual-follow", middleware.Protected(http.HandlerFunc(handlers.CheckMutualFollow))).Methods("POST")
	r.Handle("/follow", middleware.Protected(http.HandlerFunc(handlers.FollowUser))).Methods("POST")
	r.Handle("/unfollow", middleware.Protected(http.HandlerFunc(handlers.UnfollowUser))).Methods("POST")
	r.Handle("/follow-status", middleware.Protected(http.HandlerFunc(handlers.GetFollowStatus))).Methods("POST")
	r.Handle("/users", middleware.Protected(http.HandlerFunc(handlers.GetAllUsers))).Methods("GET")
	r.Handle("/users/search", middleware.Protected(http.HandlerFunc(handlers.SearchUsers))).Methods("GET")
	r.Handle("/chat-list", middleware.Protected(http.HandlerFunc(handlers.GetChatList))).Methods("GET")

	// Sync endpoints (called by auth service)
	r.HandleFunc("/sync-user", handlers.SyncUser).Methods("POST")
	r.HandleFunc("/bulk-sync-users", handlers.BulkSyncUsers).Methods("POST")

	// Start server
	log.Println("Chat service running on :8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}
