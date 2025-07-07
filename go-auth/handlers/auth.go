// package handlers

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	"go-auth/db"
// 	"go-auth/models"
// 	"go-auth/utils"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/lib/pq"
// )

// // func Signup(w http.ResponseWriter, r *http.Request) {
// // 	var user models.User
// // 	json.NewDecoder(r.Body).Decode(&user)

// // 	hashed, _ := utils.HashPassword(user.Password)
// // 	err := db.DB.QueryRow(
// // 		"INSERT INTO users(name, email, password) VALUES($1, $2, $3) RETURNING id",
// // 		user.Name, user.Email, hashed,
// // 	).Scan(&user.ID)

// // 	if err != nil {
// // 		http.Error(w, "User creation failed", http.StatusInternalServerError)
// // 		return
// // 	}

// //		json.NewEncoder(w).Encode(user)
// //	}
// func Signup(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var user models.User

// 	// Decode request body
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
// 		return
// 	}

// 	// Hash the password
// 	hashed, err := utils.HashPassword(user.Password)
// 	if err != nil {
// 		http.Error(w, `{"message":"Failed to hash password"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	// Attempt to insert user
// 	err = db.DB.QueryRow(
// 		"INSERT INTO users(name, email, password) VALUES($1, $2, $3) RETURNING id",
// 		user.Name, user.Email, hashed,
// 	).Scan(&user.ID)

// 	if err != nil {
// 		// ✅ Check for duplicate email (PostgreSQL code 23505)
// 		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" && strings.Contains(pqErr.Message, "email") {
// 			http.Error(w, `{"message":"Email already registered"}`, http.StatusConflict)
// 			return
// 		}

// 		// ✅ Fallback generic error
// 		log.Printf("Signup DB error: %v", err)
// 		http.Error(w, `{"message":"User creation failed"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	// Clear password before sending back
// 	user.Password = ""

// 	// Send back user info
// 	json.NewEncoder(w).Encode(user)
// }

// func Login(w http.ResponseWriter, r *http.Request) {
// 	var creds models.User
// 	json.NewDecoder(r.Body).Decode(&creds)

// 	var user models.User
// 	err := db.DB.QueryRow("SELECT id, password FROM users WHERE email=$1", creds.Email).Scan(&user.ID, &user.Password)
// 	if err != nil {
// 		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
// 		return
// 	}

// 	if !utils.CheckPasswordHash(creds.Password, user.Password) {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 		return
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"user_id": user.ID,
// 		"exp":     time.Now().Add(time.Hour * 72).Unix(),
// 	})

// 	tokenStr, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
// 	w.Write([]byte(tokenStr))
// }

// func GetUserInfo(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		UserID int `json:"user_id"`
// 	}
// 	json.NewDecoder(r.Body).Decode(&input)

// 	var user models.User
// 	err := db.DB.QueryRow(
// 		"SELECT name, email FROM users WHERE id = $1",
// 		input.UserID,
// 	).Scan(&user.Name, &user.Email)

// 	if err != nil {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(user)
// }

// func GetMe(w http.ResponseWriter, r *http.Request) {
// 	tokenStr := r.Header.Get("Authorization")
// 	if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
// 		http.Error(w, "Missing token", http.StatusUnauthorized)
// 		return
// 	}

// 	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
// 	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(os.Getenv("JWT_SECRET")), nil
// 	})
// 	if err != nil || !token.Valid {
// 		http.Error(w, "Invalid token", http.StatusUnauthorized)
// 		return
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
// 		return
// 	}

// 	userID, ok := claims["user_id"].(float64) // JWT stores numeric values as float64
// 	if !ok {
// 		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
// 		return
// 	}

// 	var user models.User
// 	err = db.DB.QueryRow("SELECT name, email FROM users WHERE id = $1", int(userID)).
// 		Scan(&user.Name, &user.Email)
// 	if err != nil {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	user.ID = int(userID)
// 	json.NewEncoder(w).Encode(user)
// }

// go-auth/handlers/auth.go
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go-auth/db"
	"go-auth/models"
	"go-auth/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/lib/pq"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Hash the password
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, `{"message":"Failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	// Attempt to insert user
	err = db.DB.QueryRow(
		"INSERT INTO users(name, email, password) VALUES($1, $2, $3) RETURNING id",
		user.Name, user.Email, hashed,
	).Scan(&user.ID)

	if err != nil {
		// ✅ Check for duplicate email (PostgreSQL code 23505)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" && strings.Contains(pqErr.Message, "email") {
			http.Error(w, `{"message":"Email already registered"}`, http.StatusConflict)
			return
		}

		// ✅ Fallback generic error
		log.Printf("Signup DB error: %v", err)
		http.Error(w, `{"message":"User creation failed"}`, http.StatusInternalServerError)
		return
	}

	// After successful user creation, sync to chat service
	if err := syncUserToChatService(user); err != nil {
		// Log error but don't fail registration
		log.Printf("Failed to sync user to chat service: %v", err)
	}

	// Clear password before sending back
	user.Password = ""

	// Send back user info
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	json.NewDecoder(r.Body).Decode(&creds)

	var user models.User
	err := db.DB.QueryRow("SELECT id, password FROM users WHERE email=$1", creds.Email).Scan(&user.ID, &user.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenStr, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	w.Write([]byte(tokenStr))
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID int `json:"user_id"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	err := db.DB.QueryRow(
		"SELECT name, email FROM users WHERE id = $1",
		input.UserID,
	).Scan(&user.Name, &user.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["user_id"].(float64) // JWT stores numeric values as float64
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	var user models.User
	err = db.DB.QueryRow("SELECT name, email FROM users WHERE id = $1", int(userID)).
		Scan(&user.Name, &user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.ID = int(userID)
	json.NewEncoder(w).Encode(user)
}

// Get all users endpoint
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, email FROM users ORDER BY name")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			continue
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Sync user to chat service (call this after successful registration)
func syncUserToChatService(user models.User) error {
	chatServiceURL := os.Getenv("CHAT_SERVICE_URL") + "/sync-user"

	userData := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	resp, err := http.Post(chatServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func BulkSyncUsers(w http.ResponseWriter, r *http.Request) {
	// Get all users from database
	rows, err := db.DB.Query("SELECT id, name, email FROM users ORDER BY id")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			continue
		}
		users = append(users, map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": email,
		})
	}

	// Send to chat service
	chatServiceURL := os.Getenv("CHAT_SERVICE_URL") + "/bulk-sync-users"

	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to marshal users", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(chatServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to sync users to chat service: %v", err)
		http.Error(w, "Failed to sync users", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Chat service returned non-200 status: %d", resp.StatusCode)
		http.Error(w, "Failed to sync users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("Successfully synced %d users to chat service", len(users)),
		"count":   len(users),
	})
}
