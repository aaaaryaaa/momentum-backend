// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/gorilla/mux"
// 	"github.com/joho/godotenv"

// 	"video-service/db"
// 	"video-service/handlers"
// 	"video-service/middleware"
// )

// func main() {
// 	// Load .env
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Initialize DB
// 	if err := db.Init(); err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	// Router
// 	r := mux.NewRouter()

// 	// Protected routes
// 	r.Handle("/generate-upload-url", middleware.Protected(http.HandlerFunc(handlers.GenerateUploadURL))).Methods("POST")

// 	// Protected routes
// 	r.Handle("/videos", middleware.Protected(http.HandlerFunc(handlers.GetUserVideos))).Methods("GET")

// 	// Start server
// 	port := os.Getenv("PORT")
// 	log.Println("Video service running on port", port)
// 	log.Fatal(http.ListenAndServe(":"+port, r))
// }

// main.go - Add the delete route
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"video-service/db"
	"video-service/handlers"
	"video-service/middleware"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize DB
	if err := db.Init(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Router
	r := mux.NewRouter()

	// Protected routes
	r.Handle("/generate-upload-url", middleware.Protected(http.HandlerFunc(handlers.GenerateUploadURL))).Methods("POST")
	r.Handle("/videos", middleware.Protected(http.HandlerFunc(handlers.GetUserVideos))).Methods("GET")
	r.Handle("/videos/{id}", middleware.Protected(http.HandlerFunc(handlers.DeleteVideo))).Methods("DELETE")

	// Start server
	port := os.Getenv("PORT")
	log.Println("Video service running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
