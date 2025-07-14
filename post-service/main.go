package main

import (
	"log"
	"net/http"
	"os"
	"post-service/db"
	"post-service/handlers"
	"post-service/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found")
	}

	err := db.Init()
	if err != nil {
		log.Fatal("DB init failed:", err)
	}

	r := mux.NewRouter()

	// Public endpoints
	r.Handle("/feed", middleware.Protected(http.HandlerFunc(handlers.GetFeed))).Methods("GET")

	// Protected post routes
	r.Handle("/posts", middleware.Protected(http.HandlerFunc(handlers.CreatePost))).Methods("POST")
	r.Handle("/react", middleware.Protected(http.HandlerFunc(handlers.AddReaction))).Methods("POST")
	r.Handle("/react", middleware.Protected(http.HandlerFunc(handlers.RemoveReaction))).Methods("DELETE")
	r.Handle("/nudge", middleware.Protected(http.HandlerFunc(handlers.AddNudge))).Methods("POST")
	r.Handle("/nudge", middleware.Protected(http.HandlerFunc(handlers.RemoveNudge))).Methods("DELETE")
	r.Handle("/streak", middleware.Protected(http.HandlerFunc(handlers.GetUserStreak))).Methods("GET")
	r.Handle("/leaderboard", middleware.Protected(http.HandlerFunc(handlers.GetLeaderboard))).Methods("GET")
	r.Handle("/streak-leaderboard", middleware.Protected(http.HandlerFunc(handlers.GetStreakLeaderboard))).Methods("GET")
	r.Handle("/posts/{id}", middleware.Protected(http.HandlerFunc(handlers.DeletePost))).Methods("DELETE")
	r.Handle("/thread/{id}", middleware.Protected(http.HandlerFunc(handlers.GetThread))).Methods("GET")
	r.Handle("/mythreads", middleware.Protected(http.HandlerFunc(handlers.GetMyThreads))).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Println("Post service running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
