// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/gorilla/mux"
// 	"github.com/joho/godotenv"

// 	"go-auth/db"
// 	"go-auth/handlers"
// 	"go-auth/middleware"
// )

// func main() {
// 	godotenv.Load()

// 	if err := db.Init(); err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	r := mux.NewRouter()

// 	r.HandleFunc("/signup", handlers.Signup).Methods("POST")
// 	r.HandleFunc("/login", handlers.Login).Methods("POST")
// 	r.Handle("/protected", middleware.Protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("You are authorized!"))
// 	}))).Methods("GET")

// 	log.Println("Server running on port", os.Getenv("PORT"))
// 	http.ListenAndServe(":"+os.Getenv("PORT"), r)
// }

// go-auth/main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"go-auth/db"
	"go-auth/handlers"
	"go-auth/middleware"
)

func main() {
	godotenv.Load()

	if err := db.Init(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/signup", handlers.Signup).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	r.Handle("/get-user", middleware.Protected(http.HandlerFunc(handlers.GetUserInfo))).Methods("POST")

	r.Handle("/me", middleware.Protected(http.HandlerFunc(handlers.GetMe))).Methods("GET")
	r.HandleFunc("/users", handlers.GetAllUsers).Methods("GET")
	r.HandleFunc("/bulk-sync-users", handlers.BulkSyncUsers).Methods("POST")

	r.Handle("/protected", middleware.Protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You are authorized!"))
	}))).Methods("GET")

	log.Println("Server running on port", os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
